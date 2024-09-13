import { Request, Response } from 'express';
import App from '../../app';
import { SendFundsRequest, TransferReceivedFundsRequest } from './request-types';
import GoHandler from '../go-service';
import BlockchainService from '../blockchain-service';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from '../../types';
import dotenv from 'dotenv';
import crypto from 'crypto';

dotenv.config({ path: `.env.development` });

interface RouteHandlerConfig {
    method: 'GET' | 'POST';
    path: string;
    handler: (req: Request, res: Response) => void;
}

const sendResponse = (res: Response, status: number, message: string, data?: any) => {
    res.status(status);
    const response = {
        message,
        status,
        data,
    };
    res.send(response);
};

const sendResponseOK = (res: Response, message: string, data?: any) => {
    sendResponse(res, 200, message, data);
};

const sendResponseBadRequest = (res: Response, message: string, data?: any) => {
    sendResponse(res, 400, message, data);
};

const routeHandlers = (app: App): RouteHandlerConfig[] => [
    {
        method: 'GET',
        path: '/',
        handler: (req: Request, res: Response) => {
            sendResponseOK(res, 'Service running', { timestamp: Date.now()});
        }
    },
    {
        method: 'POST',
        path: '/register-address',
        handler: async (req: Request, res: Response) => {
            const { id, K, V } = req.body;

            try {
                await app.blockchainService.registerMetaAddress(id, K, V);
                sendResponseBadRequest(res, 'Meta address registered', { id });
            } catch (err: any) {
                sendResponseBadRequest(res, err.message, { timestamp: Date.now()});
            }
        }
    },
    {
        method: 'POST',
        path: '/send',
        handler: async (req: Request, res: Response) => {
            const { 
                recipientIdType,
                id,
                recipientK,
                recipientV,
                amount,
                withProxy
            } = (req.body as SendFundsRequest);

            if (typeof amount !== 'number' || (id != null && typeof id !== 'string')) {
                return sendResponseBadRequest(res, 'Invalid request body', null);
            }

            const goHandler = app.goHandler;

            const senderRandomness = crypto.randomBytes(32).toString('hex');

            let recK, recV;

            if (recipientIdType !== 'meta_address') {
                recK = recipientK;
                recV = recipientV;
            } else {
                const resolved = await app.blockchainService.resolveMetaAddress(id!)
                const { K: resolvedK, V: resolvedV } = resolved!;

                recK = resolvedK;
                recV = resolvedV;
            }

            try {
                const sendInfo: SendInfo = await goHandler.send(senderRandomness!, recK!, recV!);

                let receipt;

                if (withProxy) {
                    receipt = await app.blockchainService.sendEthViaProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, amount.toString());
                } else {
                    receipt = await app.blockchainService.ethSentWithoutProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, amount.toString());
                }
s
                await app.db.models.sentTransactions.create({
                    transaction_hash: receipt.hash,
                    block_number: receipt.blockNumber,
                    amount: amount,
                    recipient_identifier: recipientIdType === 'meta_address' ? id : 'meta',
                    recipient_identifier_type: recipientIdType,
                    recipient_k: recipientK,
                    recipient_v: recipientV,
                    recipient_stealth_address: sendInfo.address,
                    ephemeral_key: sendInfo.pubKey,
                })


                app.loggerService.logger.info(`Sending ${amount} to stealth address: ${sendInfo.address}`);

                app.loggerService.logger.info(`Registering ephemeral key: ${sendInfo.pubKey}`);

                sendResponseOK(res, 'Transfer simulated successfully', {
                    stealthAddress: sendInfo.address,
                    ephemeralPubKey: sendInfo.pubKey,
                    viewTag: sendInfo.viewTag,
                    amount: amount
                });
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
    {
        method: 'GET',
        path: '/check-received',
        handler: async (req: Request, res: Response) => {
            const { 
                fromBlock,
                toBlock,
            } = req.query;

            app.loggerService.logger.info({fromBlock, toBlock});

            const fromBlockNumber = parseInt((fromBlock || '0') as string);
            const toBlockNumber = parseInt((toBlock || await app.blockchainService.provider.getBlockNumber()) as string);

            console.log({ fromBlockNumber, toBlockNumber });

            if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
                return sendResponseBadRequest(res, 'Invalid block numbers', null);
            }

            const goHandler = app.goHandler;

            const senderInfo = await goHandler.genSenderInfo();
            const recipientInfo = await goHandler.genRecipientInfo();

            try {
                let allReceipts: any[] = [];

                const existingReceipts = await app.db.models.receivedTransactions.findAll({
                    where: { block_number: {
                        [Op.between]: [fromBlockNumber, toBlockNumber]
                    }
                }
                });

                const newReceipt = await goHandler.receiveScan(recipientInfo.k, recipientInfo.v, [senderInfo.R], [senderInfo.viewTag]);

                const balance = await app.blockchainService.getBalance((newReceipt as any).address);
                if (balance > 0) {
                    const res = await app.db.models.sentTransactions.findAll({
                        where: {
                            recipient_stealth_address: (newReceipt as any).address,
                            amount: balance,
                        }
                    });
                    await app.db.models.receivedTransactions.create({
                        transaction_hash: res[0].transaction_hash,
                        block_number: res[0].block_number,
                        amount: balance,
                        stealth_address: (newReceipt as any).address,
                        ephemeral_key: res[0].ephemeral_key,
                        view_tag: res[0].view_tag,
                    });
                    allReceipts = [...existingReceipts, newReceipt];
                } else {
                    allReceipts = existingReceipts;
                }

                sendResponseOK(res, 'Success', { receipts: allReceipts });
            } catch (err) {
                sendResponseBadRequest(res, `Request failed: ${(err as Error).message}`, null);
            }
        }
    },
    {
        method: 'GET',
        path: '/transfer/:receiptId',
        handler: async (req: Request, res: Response) => {
            const receiptId: number = parseInt(req.params.receiptId);
            const { 
                address,
                amount,
            } = req.body as TransferReceivedFundsRequest;

            try {
                const receipt = await app.db.models.receivedTransactions.findByPk({
                    id: receiptId
                });
                if (!receipt) {
                    return sendResponseBadRequest(res, 'Receipt not found', null);
                }

                const goHandler = app.goHandler;

                const ephemeralKey = receipt.ephemeral_key;
                const viewTag = receipt.view_tag;

                const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(process.env.k!, process.env.v!, [ephemeralKey], [viewTag]);

                const tx = await app.blockchainService.transferEth(address, amount.toString(), receiveScanInfo[0].privKey);

                sendResponseOK(res, 'Success')
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
];

export default routeHandlers;