import { Server as SocketServer } from 'socket.io';
import { describe, expect, test } from '@jest/globals';
import { createServer } from 'http';
import App from '../src/app';
import configLoader from '../utils/config-loader';
import { deployContracts } from './ganache-deployment';
import { Server } from 'ganache';
import { Config } from '../types';
import { io } from 'socket.io-client';
import { Info } from '../src/types';
import SocketsHandler from '../src/services/sockets';

// Application object
let app: App;

let socketsHandler: SocketsHandler;

let ganacheServer: Server;

let config: Config;

describe('Socket.IO functionalities test', () => {
    beforeAll(async () => {
        config = configLoader.load('test');

        const {
            deployedContracts,
            privateKey,
            ganacheServer: blockchainServer,
        } = await deployContracts();

        ganacheServer = blockchainServer;

        config.blockchainConfig.privateKey = privateKey;
        config.blockchainConfig.deployedContracts = deployedContracts;

        app = new App(config);

        // Create HTTP server and Socket.IO server
        const httpServer = createServer();
        const socketServer = new SocketServer(httpServer);
        socketsHandler = new SocketsHandler(socketServer, app);

        httpServer.listen(3000); // Specify the port here

        socketsHandler.setupHandlers();

        // Start application
        await app.start();
    }, 30000);

    test('Socket connection', (done) => {
        const clientSocket = io('http://localhost:3000');

        clientSocket.on('connect', () => {
            expect(clientSocket.connected).toBe(true);
            clientSocket.disconnect();
            done();
        });
    });

    test('Check service status', (done) => {
        console.log("----------------------- SERVICE-STATUS ------------------------");
        const clientSocket = io('http://localhost:3000');

        clientSocket.emit('service-status', (response: any) => {
            console.log("service-status res", response);
            expect(response.message).toBe('Service running');
            clientSocket.disconnect();
            done();
        });
    });

    test('Register address', (done) => {
        console.log("----------------------- REGISTER ------------------------");
        const clientSocket = io('http://localhost:3000');
        
        app.goHandler.genRecipientInfo().then((recipientInfo: Info) => {
            const payload = {
                id: 'Mihailo',
                K: recipientInfo.K,
                V: recipientInfo.V,
            }
            config.stealthConfig.k = recipientInfo.k;
            config.stealthConfig.v = recipientInfo.v;

            clientSocket.emit('register-address', payload, (response: any) => {
                console.log("register-address res", response);
                expect(response.message).toBe('Meta address registered');
                clientSocket.disconnect();
                done();
            });
        }).catch((error) => {
            console.error("Error resolving Info", error);
        });
    });

    test('Send stealth transaction via Proxy', (done) => {
        console.log("----------------------- SEND ------------------------");
        const clientSocket = io('http://localhost:3000');

        const payload = {
            recipientIdType: 'id',
            id: 'Mihailo',
            amount: 10,
            withProxy: true,
        };

        clientSocket.emit('send', payload, (response: any) => {
            console.log("send res", response);
            expect(response.message).toBe('Transfer simulated successfully');
            clientSocket.disconnect();
            done();
        });
    });

    test('Check received funds', (done) => {
        console.log("----------------------- CHECK-RECEIVED ------------------------");
        const clientSocket = io('http://localhost:3000');

        clientSocket.emit('check-received', {}, (response: any) => {
            console.log("check-received res", response);
            expect(response.message).toBe('Success');
            clientSocket.disconnect();
            done();
        });
    });

    test('Transfer funds', (done) => {
        console.log("----------------------- TRANSFER ------------------------");
        const clientSocket = io('http://localhost:3000');

        const transferData = {
            receiptId: 1, // Change as needed
        };

        clientSocket.emit('transfer', transferData, (response: any) => {
            console.log("transfer res", response);
            expect(response.message).toBe('Success');
            clientSocket.disconnect();
            done();
        });
    }, 10000);

    afterAll(async () => {
        try {
            await app.stop();
            await ganacheServer.close();
        } catch (err) {
            console.error(err);
        }
    }, 50000);
});
