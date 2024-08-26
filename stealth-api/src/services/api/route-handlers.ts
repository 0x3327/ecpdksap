import { Request, Response } from 'express';
import App from '../../app';
import { SendFundsRequest, CheckReceivedRequest, TransferReceivedFundsRequest } from './request-types';

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
            res.send('Sve OK!');
        }
    },
    {
        method: 'POST',
        path: '/send',
        handler: (req: Request, res: Response) => {
            const { 
                recipientIdType,
                dns,
                address,
                recipientK,
                recipientV,
                amount,
            } = (req.body as SendFundsRequest);

            if (typeof amount !== 'number' || typeof address !== 'string') {
                return sendResponseBadRequest(res, 'Invalid request body', null);
            }

            // TODO:
            // - Generate recipient's stealth address and ephemeral key daya
            // - Send funds to stealth address
            // - Register computed ephemeral key in smart contract registry

            try {
                sendResponseOK(res, 'Transfer successfull', null)
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
    {
        method: 'GET',
        path: '/check-received',
        handler: (req: Request, res: Response) => {
            const { 
                fromBlock,
                toBlock,
            } = req.query;

            console.log(fromBlock, toBlock);

            // TODO: 
            // - Fetch from DB or Blockchain (this.app.db...)
            // - Store new receipts in db (received_transactions)

            const receivedData = {};

            try {
                sendResponseOK(res, 'Success', receivedData)
            } catch (err) {
                sendResponseBadRequest(res, `Request failed: ${(err as Error).message}`, null);
            }
        }
    },
    {
        method: 'GET',
        path: '/transfer/:receiptId',
        handler: (req: Request, res: Response) => {
            const receiptId: number = parseInt(req.params.receiptId);
            const { 
                address,
                amount,
            } = req.body as TransferReceivedFundsRequest;

            console.log(address, amount);

            // TODO: 
            // - Fetch receipt data from received_transactions
            // - Generate private key for address
            // - Send <amount> to <address>

            try {
                sendResponseOK(res, 'Success')
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
];

export default routeHandlers;