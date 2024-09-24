import { Server } from 'socket.io';
import App from '../../app';
import { SendFundsRequest, TransferReceivedFundsRequest } from './request-types';
import { Op } from 'sequelize';
import { ReceiveScanInfo, SendInfo } from '../../types';
import dotenv from 'dotenv';
import configLoader from '../../../utils/config-loader';

dotenv.config({ path: `.env.development` });

const config = configLoader.load('test');

const setupSocketHandlers = (io: Server, app: App) => {
    io.on('connection', (socket) => {
        console.log('Client connected:', socket.id);

        socket.on('service-status', (callback) => {
            callback({ message: 'Service running', timestamp: Date.now() });
        });

        socket.on('register-address', async (data: { id: string; K: string; V: string }, callback) => {
            try {
                await app.blockchainService.registerMetaAddress(data.id, data.K, data.V);
                callback({ message: 'Meta address registered', id: data.id });
            } catch (err: any) {
                callback({ error: err.message, timestamp: Date.now() });
            }
        });

        socket.on('send', async (data: SendFundsRequest, callback) => {
            const { recipientIdType, id, recipientK, recipientV, amount, withProxy } = data;

            if (typeof amount !== 'number' || (id != null && typeof id !== 'string')) {
                return callback({ error: 'Invalid request body' });
            }

            const goHandler = app.goHandler;

            let recK, recV;

            if (recipientIdType === 'meta_address') {
                recK = recipientK;
                recV = recipientV;
            } else {
                const resolved = await app.blockchainService.resolveMetaAddress(id!);
                const { K: resolvedK, V: resolvedV } = resolved!;
                recK = resolvedK;
                recV = resolvedV;
            }

            try {
                const senderInfo = await goHandler.genSenderInfo();
                const sendInfo: SendInfo = await goHandler.send(senderInfo.r, recK!, recV!);

                let receipt;
                if (withProxy) {
                    receipt = await app.blockchainService.sendEthViaProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount.toString());
                } else {
                    receipt = await app.blockchainService.ethSentWithoutProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount.toString());
                }

                await app.db.models.sentTransactions.create({
                    transaction_hash: receipt.hash,
                    block_number: receipt.blockNumber,
                    amount: amount,
                    recipient_identifier: recipientIdType === 'meta_address' ? id : 'meta',
                    recipient_identifier_type: recipientIdType,
                    recipient_k: recK,
                    recipient_v: recV,
                    recipient_stealth_address: sendInfo.address,
                    ephemeral_key: sendInfo.pubKey,
                });

                app.loggerService.logger.info(`Sending ${amount} to stealth address: ${sendInfo.address}`);
                app.loggerService.logger.info(`Registering ephemeral key: ${sendInfo.pubKey}`);

                callback({ message: 'Transfer simulated successfully', data: {
                    stealthAddress: sendInfo.address,
                    ephemeralPubKey: sendInfo.pubKey,
                    viewTag: sendInfo.viewTag,
                    amount: amount
                }});
            } catch (err) {
                callback({ error: `Transfer failed: ${(err as Error).message}` });
            }
        });

        socket.on('check-received', async (data: { fromBlock?: number; toBlock?: number }, callback) => {
            const { fromBlock, toBlock } = data;

            app.loggerService.logger.info({ fromBlock, toBlock });

            const fromBlockNumber = fromBlock || 0;
            const toBlockNumber = toBlock || await app.blockchainService.provider.getBlockNumber();

            if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
                return callback({ error: 'Invalid block numbers' });
            }

            try {
                const receipts = await app.db.models.receivedTransactions.findAll({
                    where: {
                        block_number: {
                            [Op.between]: [fromBlockNumber, toBlockNumber]
                        }
                    }
                });

                callback({ message: 'Success', receipts });
            } catch (err) {
                callback({ error: `Request failed: ${(err as Error).message}` });
            }
        });

        socket.on('transfer', async (data: { receiptId: number; address?: string; amount?: number }, callback) => {
            const receiptId: number = data.receiptId;
            const { address, amount } = data;

            try {
                const receipt = await app.db.models.receivedTransactions.findByPk(receiptId);
                if (!receipt) {
                    return callback({ error: 'Receipt not found' });
                }

                const goHandler = app.goHandler;
                const k = app.config.stealthConfig.k;
                const v = app.config.stealthConfig.v;

                const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(k, v, [receipt.ephemeral_key], [receipt.view_tag]);

                const transferAddress = address || config.stealthConfig.transferAddress;
                const transferAmount = amount || 0.001;
                const tx = await app.blockchainService.transferEth(transferAddress, transferAmount.toString(), receiveScanInfo[0].privKey);

                callback({ message: 'Success' , tx});
            } catch (err) {
                callback({ error: `Transfer failed: ${(err as Error).message}` });
            }
        });
    });
};

export default setupSocketHandlers;