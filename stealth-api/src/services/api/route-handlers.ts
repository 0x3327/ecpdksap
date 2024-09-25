import { Request, Response } from 'express';
import App from '../../app';
import { CheckReceivedRequest, RegisterAddressRequest, SendFundsRequest, TransferReceivedFundsRequest } from './request-types';
import { Op } from 'sequelize';
import { ReceiveScanInfo, SendInfo } from '../../types';
import dotenv from 'dotenv';
import configLoader from '../../../utils/config-loader';

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

const config = configLoader.load('test');

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
            const { id, K, V } = (req.body as RegisterAddressRequest);
            try {
                await app.blockchainService.registerMetaAddress(id, K, V);
                sendResponseOK(res, 'Meta address registered', { id });
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

            let recK, recV;

            if (recipientIdType === 'meta_address') {
                recK = recipientK;
                recV = recipientV;
            } else {
                const resolved = await app.blockchainService.resolveMetaAddress(id!)
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

                console.log(sendResponseOK(res, 'Transfer simulated successfully', {
                    stealthAddress: sendInfo.address,
                    ephemeralPubKey: sendInfo.pubKey,
                    viewTag: sendInfo.viewTag,
                    amount: amount
                }));
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

            if (isNaN(fromBlockNumber) || isNaN(toBlockNumber)) {
                return sendResponseBadRequest(res, 'Invalid block numbers', null);
            }

            try {
                const receipts = await app.db.models.receivedTransactions.findAll({
                    where: { 
                        block_number: {
                            [Op.between]: [fromBlockNumber, toBlockNumber]
                        }
                    }
                });

                sendResponseOK(res, 'Success', { receipts });
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
                const receipt = await app.db.models.receivedTransactions.findByPk(receiptId);
                if (!receipt) {
                    return sendResponseBadRequest(res, 'Receipt not found', null);
                }

                const goHandler = app.goHandler;
                const k = app.config.stealthConfig.k;
                const v = app.config.stealthConfig.v;

                const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(k, v, [receipt.ephemeral_key], [receipt.view_tag]);
                
                const transferAddress = address || config.stealthConfig.transferAddress;
                const transferAmount = amount || 0.001;
                const tx = await app.blockchainService.transferEth(transferAddress, transferAmount.toString(), receiveScanInfo[0].privKey);

                console.log("tx", tx);

                sendResponseOK(res, 'Success')
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
];

export default routeHandlers;