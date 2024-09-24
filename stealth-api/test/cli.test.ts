import { describe, expect, test } from '@jest/globals';
import App from '../src/app';
import configLoader from '../utils/config-loader';
import { deployContracts } from './ganache-deployment';
import { Server } from 'ganache';
import { Config } from '../types';
import { program, registerCommands } from '../src/services/api/cli';

// Application object
let app: App;

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
        registerCommands(app);

        // Start application
        await app.start();

    }, 30000)

    test('Starting application', () => {
       expect(app).not.toBeNull();
    });

    test('Check heartbeat command', async () => {
        try {
            console.log("----------------------- SERVICE-STATUS ------------------------");
            await program.parseAsync(['node', 'test', 'service-status']);
            // expect(data.message).toBe('Service running');
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

            await program.parseAsync(['node', 'test', 'register-address', '--id', payload.id, '--K', payload.K, '--V', payload.V]);

            // Wait for MetaAddressRegistry event
            await (new Promise((resolve, reject) => setTimeout(resolve, 5000)));
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }
    }, 10000);

    test('Send stealth transaction via Proxy', async () => {
        console.log("----------------------- SEND ------------------------");
        try {
            const payload = {
                recipientIdType: 'id',
                id: 'Mihailo', 
                amount: '10',
                withProxy: true,
            }

            await program.parseAsync(['node', 'test', 'send', '--recipientIdType', payload.recipientIdType,
                                    '--id', payload.id, '--amount', payload.amount, '--withProxy']);

            // Wait for Announcement event
            await (new Promise((resolve, reject) => setTimeout(resolve, 5000)));

        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }   
    }, 10000);

    test('Check received funds', async () => {
        console.log("----------------------- CHECK-RECEIVED ------------------------");
        try {
            await program.parseAsync(['node', 'test', 'check-received'])
        } catch(err) {
            console.log(err);
            expect(true).toBe(false);
        }
    });

    test('Transfer funds', async () => {
        console.log("----------------------- TRANSFER ------------------------");
        try {
            await program.parseAsync(['node', 'test', 'transfer', '--receiptId', '1']);
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }
    });

    afterAll(async () => {
        try {
            await app.blockchainService.provider.destroy();
            await app.stop()
            await ganacheServer.close();
        } catch (err) {}
    }, 50000)
});
