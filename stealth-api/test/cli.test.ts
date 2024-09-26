import { describe, expect, test } from '@jest/globals';
import App from '../src/app';
import configLoader from '../utils/config-loader';
import { deployContracts } from './ganache-deployment';
import { Server } from 'ganache';
import { Config } from '../types';
import CLIService from '../src/services/cli';

// Application object
let app: App;

let cliService : CLIService;

let ganacheServer: Server

let config: Config;

describe('CLI API commands test', () => {
    beforeAll(async () => {
        config = configLoader.load('test');
        
        const {
            deployedContracts,
            privateKey,
            ganacheServer: blockchainServer,
        } = await deployContracts();

        // Ganache server
        ganacheServer = blockchainServer;

        config.blockchainConfig.privateKey = privateKey;
        config.blockchainConfig.deployedContracts = deployedContracts;

        app = new App(config);
        cliService = new CLIService(app);

        // Start application
        await app.start();

    }, 30000)

    test('Starting application', () => {
       expect(app).not.toBeNull();
    });

    test('Check heartbeat command', async () => {
        try {
            console.log("----------------------- SERVICE-STATUS ------------------------");
            process.argv = ['node', 'test', 'service-status'];
            await cliService.serviceStatus();
            await cliService.run();
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }    
    });

    test('Register', async () => {
        console.log("----------------------- REGISTER ------------------------");
        try {
            const recipientInfo = await app.goHandler.genRecipientInfo();
            config.stealthConfig.k = recipientInfo.k;
            config.stealthConfig.v = recipientInfo.v;

            const payload = {
                id: 'Mihailo',
                K: recipientInfo.K,
                V: recipientInfo.V,
            };

            process.argv = ['node', 'test', 'register-address', '--id', payload.id, '--K', payload.K, '--V', payload.V];
            await cliService.registerAddress();
            await cliService.run();

            // Wait for MetaAddressRegistry event
            await (new Promise((resolve, reject) => setTimeout(resolve, 5000)));
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }
    }, 6000);

    test('Send stealth transaction via Proxy', async () => {
        console.log("----------------------- SEND ------------------------");
        try {
            const payload = {
                recipientIdType: 'id',
                id: 'Mihailo', 
                amount: '10',
                withProxy: true,
            }

            process.argv = ['node', 'test', 'send', '--recipientIdType', payload.recipientIdType,
                                    '--id', payload.id, '--amount', payload.amount, '--withProxy'];
            await cliService.sendFunds()
            await cliService.run();

            // Wait for Announcement event
            await (new Promise((resolve, reject) => setTimeout(resolve, 5000)));

        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }   
    }, 6000);

    test('Check received funds', async () => {
        console.log("----------------------- CHECK-RECEIVED ------------------------");
        try {
            const fromBlock = 0;
            const toBlock = 5;
            process.argv =  ['node', 'test', 'check-received'/*, '--fromBlock', fromBlock.toString(), '--toBlock', toBlock.toString()*/];
            await cliService.checkReceived();
            await cliService.run();
        } catch(err) {
            console.log(err);
            expect(true).toBe(false);
        }
    });

    test('Transfer funds', async () => {
        console.log("----------------------- TRANSFER ------------------------");
        try {
            process.argv = ['node', 'test', 'transfer', '--receiptId', '1'];
            await cliService.transfer();
            await cliService.run();
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }
    });

    afterAll(async () => {
        try {
            await app.stop();
            await ganacheServer.close();
        } catch (err) {}
    }, 50000)
});
