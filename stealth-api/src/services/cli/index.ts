import { Command } from "commander";
import App from "../../app";
import { Info, ReceiveScanInfo, SendInfo } from "../../types";
import { Op } from "sequelize";

const sendResponse = (status: number, message: string, data?: any) => {
    return { status, message, data };
};

const sendResponseOK = (message: string, data?: any) => {
    return sendResponse(200, message, data);
};

const sendResponseBadRequest = (message: string, data?: any) => {
    return sendResponse(400, message, data);
};

class CLIService {
    private app: App;
    private program: Command;

    constructor(app: App) {
        this.app = app;
        this.program = new Command;
    }

    public async serviceStatus() {
        this.program
            .command('service-status')
            .description('Check if the service is running')
            .action(() => {
                const response = sendResponseOK('Service running', { timestamp: Date.now() });
                console.log(response);
            });
    }

    public async registerAddress() {
        this.program
            .command('register-address')
            .description('Register a meta address')
            .requiredOption('--id <string>', 'ID')
            .requiredOption('--K <string>', 'Spending public key')
            .requiredOption('--V <string>', 'Viewwing public key')
            .action(async (opts) => {
                try {
                    const { id, K, V } = opts;
                    await this.app.blockchainService.registerMetaAddress(id, K, V);
                    const response = sendResponseOK('Meta address registered', { id });
                    console.log(response);
                } catch (err: any) {
                    const response = sendResponseBadRequest(err.message, { timestamp: Date.now() });
                    console.log(response);
                }
            });
    }

    public async sendFunds() {
        this.program
            .command('send')
            .description('Send funds')
            .requiredOption('--recipientIdType <string>', 'Recipient ID Type')
            .option('--id <string>', 'ID')
            .option('--recipientK <string>', 'Recipient spending public key')
            .option('--recipientV <string>', 'Recipient viewing public key')
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

                const goHandler = this.app.goHandler;
                let recK, recV;

                if (recipientIdType === 'meta_address') {
                    recK = recipientK;
                    recV = recipientV;
                } else {
                    const resolved = await this.app.blockchainService.resolveMetaAddress(id);
                    recK = resolved!.K;
                    recV = resolved!.V;
                }

                try {
                    const senderInfo: Info = await goHandler.genSenderInfo();
                    const sendInfo: SendInfo = await goHandler.send(senderInfo.r, recK, recV);

                    let receipt;
                    if (withProxy) {
                        receipt = await this.app.blockchainService.sendEthViaProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount);
                    } else {
                        receipt = await this.app.blockchainService.ethSentWithoutProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount);
                    }

                    await this.app.db.models.sentTransactions.create({
                        transaction_hash: receipt.hash,
                        block_number: receipt.blockNumber,
                        amount: amount,
                        recipient_identifier: recipientIdType === 'meta_address' ? id : 'meta',
                        recipient_identifier_type: recipientIdType,
                        recipient_k: recK,
                        recipient_v: recV,
                        recipient_stealth_address: sendInfo.address,
                        ephemeral_key: sendInfo.pubKey,
                    })
    
                    this.app.loggerService.logger.info(`Sending ${amount} to stealth address: ${sendInfo.address}`);
    
                    this.app.loggerService.logger.info(`Registering ephemeral key: ${sendInfo.pubKey}`);

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
    }

    public async checkReceived() {
        this.program
            .command('check-received')
            .description('Check received transactions')
            .option('--fromBlock <string>', 'From block number')
            .option('--toBlock <string>', 'To block number')
            .action(async (opts) => {
                const { fromBlock, toBlock } = opts;

                const fromBlockNumber = parseInt(fromBlock || '0');
                const toBlockNumber = parseInt(toBlock || await this.app.blockchainService.provider.getBlockNumber());

                if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
                    console.log(sendResponseBadRequest('Invalid block numbers', null));
                    return;
                }

                try {
                    const receipts = await this.app.db.models.receivedTransactions.findAll({
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
    }

    public async transfer() {
        this.program
            .command('transfer')
            .description('Transfer received funds')
            .requiredOption('--receiptId <number>', 'Receipt ID')
            .option('--address <string>', 'Transfer address')
            .option('--amount <number>', 'Transfer amount')
            .action(async (opts) => {
                const { receiptId, address, amount } = opts;

                try {
                    const receipt = await this.app.db.models.receivedTransactions.findByPk(receiptId);
                    if (!receipt) {
                        console.log(sendResponseBadRequest('Receipt not found', null));
                        return;
                    }

                    const goHandler = this.app.goHandler;
                    const k = this.app.config.stealthConfig.k;
                    const v = this.app.config.stealthConfig.v;

                    const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(k, v, [receipt.ephemeral_key], [receipt.view_tag]);
                    const transferAddress = address || this.app.config.stealthConfig.transferAddress;
                    const transferAmount = amount || 0.001;

                    const tx = await this.app.blockchainService.transferEth(transferAddress, transferAmount.toString(), receiveScanInfo[0].privKey);

                    console.log(sendResponseOK('Success', { tx }));
                } catch (err) {
                    console.log(sendResponseBadRequest(`Transfer failed: ${(err as Error).message}`, null));
                }
            });
    }

    public run() {
        this.program.parse(process.argv);
    }
}

export default CLIService;