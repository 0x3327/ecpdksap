import { Server as SocketServer } from 'socket.io';
import { createServer, Server } from 'http';
import App from '../../app';
import { SendFundsRequest } from '../api/request-types';
import { Op } from 'sequelize';
import { ReceiveScanInfo, SendInfo } from '../../types';

class SocketsService {
    private io: SocketServer;
    private app: App;

    constructor(app: App) {
        const httpServer = createServer();
        this.io = new SocketServer(httpServer);
        this.app = app;
    }

    public start() {
        const { port } = this.app.config.socketConfig;
        this.io.listen(port);
        this.app.loggerService.logger.info(`Socket service listening on ${port}`)

        this.io.on('connection', (socket) => {
            console.log('Client connected:', socket.id);
            socket.on('service-status', (callback) => this.serviceStatus(socket, callback));
            socket.on('register-address', (data, callback) => this.registerAddress(socket, data, callback));
            socket.on('send', (data, callback) => this.sendFunds(socket, data, callback));
            socket.on('check-received', (data, callback) => this.checkReceived(socket, data, callback));
        });
    }

    private serviceStatus(socket: any, callback: Function) {
        callback({ message: 'Service running', timestamp: Date.now() });
    }

    private async registerAddress(socket: any, data: { id: string; K: string; V: string }, callback: Function) {
        console.log(data, callback);
        try {
            await this.app.blockchainService.registerMetaAddress(data.id, data.K, data.V);
            callback({ message: 'Meta address registered', id: data.id });
        } catch (err: any) {
            callback({ error: err.message, timestamp: Date.now() });
        }
    }

    private async sendFunds(socket: any, data: SendFundsRequest, callback: Function) {
        const { recipientIdType, id, recipientK, recipientV, amount, withProxy } = data;

        if (typeof amount !== 'number' || (id != null && typeof id !== 'string')) {
            return callback({ error: 'Invalid request body' });
        }

        const goHandler = this.app.goHandler;

        let recK, recV;

        if (recipientIdType === 'meta_address') {
            recK = recipientK;
            recV = recipientV;
        } else {
            const resolved = await this.app.blockchainService.resolveMetaAddress(id!);
            const { K: resolvedK, V: resolvedV } = resolved!;
            recK = resolvedK;
            recV = resolvedV;
        }

        try {
            const senderInfo = await goHandler.genSenderInfo();
            const sendInfo: SendInfo = await goHandler.send(senderInfo.r, recK!, recV!);

            let receipt;
            if (withProxy) {
                receipt = await this.app.blockchainService.sendEthViaProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount.toString());
            } else {
                receipt = await this.app.blockchainService.ethSentWithoutProxy(sendInfo.address, senderInfo.R, sendInfo.viewTag, amount.toString());
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
            });

            this.app.loggerService.logger.info(`Sending ${amount} to stealth address: ${sendInfo.address}`);

            this.app.loggerService.logger.info(`Registering ephemeral key: ${sendInfo.pubKey}`);

            callback({
                message: 'Transfer simulated successfully', data: {
                    stealthAddress: sendInfo.address,
                    ephemeralPubKey: sendInfo.pubKey,
                    viewTag: sendInfo.viewTag,
                    amount: amount
                }
            });
        } catch (err) {
            callback({ error: `Transfer failed: ${(err as Error).message}` });
        }
    }

    private async checkReceived(socket: any, data: { fromBlock?: number; toBlock?: number }, callback: Function) {
        const { fromBlock, toBlock } = data;

        this.app.loggerService.logger.info({ fromBlock, toBlock });

        const fromBlockNumber = fromBlock || 0;
        const toBlockNumber = toBlock || await this.app.blockchainService.provider.getBlockNumber();

        if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
            return callback({ error: 'Invalid block numbers' });
        }

        try {
            const receipts = await this.app.db.models.receivedTransactions.findAll({
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
    }

    private async transfer(socket: any, data: { receiptId: number; address?: string; amount?: number }, callback: Function) {
        const receiptId: number = data.receiptId;
        const { address, amount } = data;

        try {
            const receipt = await this.app.db.models.receivedTransactions.findByPk(receiptId);
            if (!receipt) {
                return callback({ error: 'Receipt not found' });
            }

            const goHandler = this.app.goHandler;
            const k = this.app.config.stealthConfig.k;
            const v = this.app.config.stealthConfig.v;

            const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(k, v, [receipt.ephemeral_key], [receipt.view_tag]);

            const transferAddress = address || this.app.config.stealthConfig.transferAddress;
            const transferAmount = amount || 0.001;
            const tx = await this.app.blockchainService.transferEth(transferAddress, transferAmount.toString(), receiveScanInfo[0].privKey);

            callback({ message: 'Success', tx });
        } catch (err) {
            callback({ error: `Transfer failed: ${(err as Error).message}` });
        }
    }

    public stop(): Promise<void> {
        return new Promise((resolve, reject) => {
            this.io.httpServer.close((err) => {
                if (err) {
                    reject(err);
                } else {
                    resolve();
                }
            });
        })
    }
}

export default SocketsService;