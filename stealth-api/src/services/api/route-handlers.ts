import { Request, Response } from 'express';
import App from '../../app';
import { SendFundsRequest, TransferReceivedFundsRequest } from './request-types';
import GoHandler from '../go-service';
import BlockchainService from '../blockchain-service';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from '../../types'; 
import dotenv from 'dotenv';
import crypto from 'crypto';
import configLoader from '../../../utils/config-loader';
import { mulPointEscalar, Base8, Point } from '@zk-kit/baby-jubjub'; 

 
const generateKeys = (): { privateKey: bigint, publicKey: Point<bigint> } => { 
    const privateKey = BigInt('0x' + crypto.randomBytes(32).toString('hex')) % 
    BigInt("2736030358979909402780800718157159386076813972158567259200215660948447373041"); // Generate random private key 
    const publicKey = mulPointEscalar(Base8, privateKey); // Compute public key by scalar multiplication 
    return { privateKey, publicKey }; 
}; 
 
const { privateKey, publicKey } = generateKeys(); 
const publicKeyX = publicKey[0].toString(10); 
const publicKeyY = publicKey[1].toString(10); 
 
console.log("Generated private key:", privateKey.toString(10)); 
console.log("Generated public key:", publicKeyX, publicKeyY);

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
        path: '/send',
        handler: async (req: Request, res: Response) => {
            const { 
                recipientIdType,
                ens,
                address,
                recipientK,
                recipientV,
                amount,
                withProxy
            } = (req.body as SendFundsRequest);

            if (typeof amount !== 'number' || (address != null && typeof address !== 'string') || (ens != null && typeof ens !== 'string')) {
                return sendResponseBadRequest(res, 'Invalid request body', null);
            }

            const goHandler = app.goHandler;

            const senderInfo = await goHandler.genSenderInfo();
            config.stealthConfig.senderRandomness = senderInfo.r;
            config.stealthConfig.senderR = senderInfo.R;

            try {
                const sendInfo: SendInfo = await goHandler.send(config.stealthConfig.senderRandomness!, config.stealthConfig.recipientK!, config.stealthConfig.recipientV!);

                config.stealthConfig.Rs = [config.stealthConfig.senderR];
                config.stealthConfig.ViewTags = [sendInfo.viewTag]

                let receipt;

                if (withProxy) {
                    receipt = await app.blockchainService.sendEthViaProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, amount.toString());
                } else {
                    receipt = await app.blockchainService.ethSentWithoutProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, amount.toString());
                }

                await app.db.models.sentTransactions.create({
                    transaction_hash: receipt.hash,
                    block_number: receipt.blockNumber,
                    amount: amount,
                    recipient_identifier: recipientIdType === 'eth_ens' ? ens : address,
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

            try {
                let allReceipts: any[] = [];

                const existingReceipts = await app.db.models.receivedTransactions.findAll({
                    where: { 
                        block_number: {
                            [Op.between]: [fromBlockNumber, toBlockNumber]
                    }
                }
                });

                const newReceipts = await goHandler.receiveScan(config.stealthConfig.recipientk, config.stealthConfig.recipientv, config.stealthConfig.Rs, config.stealthConfig.ViewTags);

                // TODO: add for loop for checking every element in newReceipt array
                const balance = await app.blockchainService.getBalance((newReceipts[0] as any).address);
                if (balance > 0) {
                    const res = await app.db.models.sentTransactions.findAll({
                        where: {
                            recipient_stealth_address: (newReceipts[0] as any).address,
                            // amount: balance,
                        }
                    });
                    // console.log("nasao u db rec tx", res_rec);
                    // await app.db.models.receivedTransactions.create({
                    //     transaction_hash: res[0].transaction_hash,
                    //     block_number: res[0].block_number,
                    //     amount: balance,
                    //     stealth_address: (newReceipts[0] as any).address,
                    //     ephemeral_key: res[0].ephemeral_key,
                    //     view_tag: config.stealthConfig.ViewTags[0],
                    // });
                    // console.log("dodao u db received tx");
                    allReceipts = [...existingReceipts, newReceipts[0]];
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
                const receipt = await app.db.models.receivedTransactions.findByPk(receiptId);
                if (!receipt) {
                    return sendResponseBadRequest(res, 'Receipt not found', null);
                }

                const goHandler = app.goHandler;

                const ephemeralKeyHex = receipt.ephemeral_key;
                console.log("ephemeralKey", ephemeralKeyHex);
                console.log("Rs", config.stealthConfig.Rs);
                const ephemeralKeySliced = ephemeralKeyHex.slice(2);
                const ephemeralKey = ephemeralKeySliced.replace('e', '.');
                const viewTagHex = receipt.view_tag;
                console.log("viewTag", viewTagHex);
                console.log("ViewTags", config.stealthConfig.ViewTags);
                const viewTag = viewTagHex.slice(2);
                const receiveScanInfo: ReceiveScanInfo[] = await goHandler.receiveScan(config.stealthConfig.recipientk!, config.stealthConfig.recipientv!, [config.stealthConfig.senderR], [viewTag]);
                const addressDefined = address || config.stealthConfig.transferAddress;
                const amountDefined = amount || 10;
                console.log("address", addressDefined);
                console.log("amount", amountDefined);
                console.log("privKey", receiveScanInfo[0].privKey);
                const tx = await app.blockchainService.transferEth(addressDefined, amountDefined.toString(), receiveScanInfo[0].privKey);

                console.log("tx", tx);

                sendResponseOK(res, 'Success')
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
    {
        method: 'POST',
        path: '/register-account',
        handler: async (req: Request, res: Response) => {
            try {
                // Generate private and public keys for the user
                const { privateKey, publicKey } = generateKeys();
                const publicKeyX = publicKey[0];
                const publicKeyY = publicKey[1];

                //console.log("Generated private key:", privateKey.toString(16));
                //console.log("Generated public key:", publicKeyX, publicKeyY);

                // Generate necessary data
                const data = {
                    name: req.body.name,
                    pid: req.body.pid, 
                    publicKeyX: publicKeyX.toString(),
                    publicKeyY: publicKeyY.toString()
                };

                // Forward the request to regulatory body with the necessary data
                // TODO
                const response = await fetch("http://localhost:5555/register-user", {
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(data),
                    method: "POST"
                })
                const responseData = await response.json();
                const proofData = responseData["data"];
                sendResponseOK(res, "Recieved merkle proof", proofData);

            }catch(error){
                console.error("Error during registration:", error);
                res.status(500).json({ error: "Registration failed" });
            }
        }
    },
];

export default routeHandlers;