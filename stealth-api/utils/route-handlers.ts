import { Request, Response } from 'express';
import App from '../src/app';

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
            const { amount, address } = req.body;

            if (typeof amount !== 'number' || typeof address !== 'string') {
                return sendResponseBadRequest(res, 'Invalid request body', null);
            }

            try {
                // app.stealthService?.send(address, amount);
                sendResponseOK(res, 'Transfer successfull', null)
            } catch (err) {
                sendResponseBadRequest(res, `Transfer failed: ${(err as Error).message}`, null);
            }
        }
    },
];

export default routeHandlers;