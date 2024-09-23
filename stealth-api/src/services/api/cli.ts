import { Command } from 'commander';
import App from '../../app';
import { Op } from 'sequelize';
import { ReceiveScanInfo, SendInfo } from '../../types';
import dotenv from 'dotenv';
import configLoader from '../../../utils/config-loader';

dotenv.config({ path: `.env.development` });

const program = new Command();
const config = configLoader.load('test');

// Helper functions for sending responses
const sendResponse = (status: number, message: string, data?: any) => {
    return { status, message, data };
};

const sendResponseOK = (message: string, data?: any) => {
    return sendResponse(200, message, data);
};

const sendResponseBadRequest = (message: string, data?: any) => {
    return sendResponse(400, message, data);
};

// Register commands
const registerCommands = (app: App) => {
    program
        .command('service-status')
        .description('Check if the service is running')
        .action(() => {
            const response = sendResponseOK('Service running', { timestamp: Date.now() });
            console.log(response);
        });

    program
        .command('register-address')
        .description('Register a meta address')
        .requiredOption('--id <string>', 'ID')
        .requiredOption('--K <string>', 'K value')
        .requiredOption('--V <string>', 'V value')
        .action(async (opts) => {
            try {
                const { id, K, V } = opts;

                await app.blockchainService.registerMetaAddress(id, K, V);
                const response = sendResponseOK('Meta address registered', { id });
                console.log(response);
            } catch (err: any) {
                const response = sendResponseBadRequest(err.message, { timestamp: Date.now() });
                console.log(response);
            }
        });

    program
        .command('send')
        .description('Send funds')
        .requiredOption('--recipientIdType <string>', 'Recipient ID Type')
        .option('--id <string>', 'ID')
        .option('--recipientK <string>', 'Recipient K')
        .option('--recipientV <string>', 'Recipient V')
        .requiredOption('--amount <string>', 'Amount to send')
        .requiredOption('--withProxy', 'Send via proxy')
        .action(async (opts) => {
            const {
                recipientIdType,
                id,
                recipientK,
                recipientV,
                amount,
                withProxy
            } = opts;

            if (typeof amount !== 'string' || (id != null && typeof id !== 'string')) {
                console.log(sendResponseBadRequest('Invalid request body', null));
                return;
            }

            const goHandler = app.goHandler;
            let recK, recV;

            if (recipientIdType === 'meta_address') {
                recK = recipientK;
                recV = recipientV;
            } else {
                const resolved = await app.blockchainService.resolveMetaAddress(id);
                recK = resolved!.K;
                recV = resolved!.V;
            }

            try {
                const senderInfo = await goHandler.genSenderInfo();
                const sendInfo: SendInfo = await goHandler.send(senderInfo.r, recK, recV);

                let receipt;
                if (withProxy) {
                    receipt = await app.blockchainService.sendEthViaProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount);
                } else {
                    receipt = await app.blockchainService.ethSentWithoutProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount);
                }

                console.log(sendResponseOK('Transfer simulated successfully', {
                    stealthAddress: sendInfo.address,
                    ephemeralPubKey: sendInfo.pubKey,
                    viewTag: sendInfo.viewTag,
                    amount: amount
                }));
            } catch (err) {
                console.log(sendResponseBadRequest(`Transfer failed: ${(err as Error).message}`, null));
            }
        });

    program
        .command('check-received')
        .description('Check received transactions')
        .option('--fromBlock <number>', 'From block number')
        .option('--toBlock <number>', 'To block number')
        .action(async (opts) => {
            const { fromBlock, toBlock } = opts;

            const fromBlockNumber = parseInt(fromBlock || '0');
            const toBlockNumber = parseInt(toBlock || await app.blockchainService.provider.getBlockNumber());

            if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
                console.log(sendResponseBadRequest('Invalid block numbers', null));
                return;
            }

            try {
                const receipts = await app.db.models.receivedTransactions.findAll({
                    where: {
                        block_number: {
                            [Op.between]: [fromBlockNumber, toBlockNumber]
                        }
                    }
                });
                console.log(sendResponseOK('Success', { receipts }));
            } catch (err) {
                console.log(sendResponseBadRequest(`Request failed: ${(err as Error).message}`, null));
            }
        });

    program
        .command('transfer')
        .description('Transfer received funds')
        .requiredOption('--receiptId <number>', 'Receipt ID')
        .option('--address <string>', 'Transfer address')
        .option('--amount <number>', 'Transfer amount')
        .action(async (opts) => {
            const { receiptId, address, amount } = opts;

            try {
                const receipt = await app.db.models.receivedTransactions.findByPk(receiptId);
                if (!receipt) {
                    console.log(sendResponseBadRequest('Receipt not found', null));
                    return;
                }

                const goHandler = app.goHandler;
                const k = app.config.stealthConfig.k;
                const v = app.config.stealthConfig.v;

                const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(k, v, [receipt.ephemeral_key], [receipt.view_tag]);
                const transferAddress = address || config.stealthConfig.transferAddress;
                const transferAmount = amount || 0.001;

                const tx = await app.blockchainService.transferEth(transferAddress, transferAmount.toString(), receiveScanInfo[0].privKey);

                console.log("tx", tx);

                console.log(sendResponseOK('Success'));
            } catch (err) {
                console.log(sendResponseBadRequest(`Transfer failed: ${(err as Error).message}`, null));
            }
        });
    }

export { program, registerCommands };

// Add this condition to prevent parsing when imported
if (require.main === module) {
    program.parse(process.argv);
}