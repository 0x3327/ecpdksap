import { describe, expect, test } from '@jest/globals';
import App from '../src/app';
import { after, before } from 'node:test';
import configLoader from '../utils/config-loader';
import axios, { AxiosInstance } from 'axios';
import { deployContracts } from './ganache-deployment';
import { Server } from 'ganache';

// Application object
let app: App;

// Axios instance
let axiosInstance: AxiosInstance;

let ganacheServer: Server

describe('API routes test', () => {
    beforeAll(async () => {
        const config = configLoader.load('test');
        
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

        // Start application
        await app.start();

        // Read API config params
        const { host, port } = app.config.apiConfig;

        // Initialize axios instance
        axiosInstance = axios.create({
            baseURL: `http://${host}:${port}`
          });

        // console.log(axiosInstance);
    })

    test('Starting application', () => {
       expect(app).not.toBeNull();
    });

    test('Check heartbeat route', async () => {
        try {
            const res = await axiosInstance.get('/');
            const { data } = res;

            expect(data.message).toBe('Service running');
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }    
    })

    test('Send stealth transaction', async () => {
        try {
            const recipientInfo = await app.goHandler.genRecipientInfo();
            
            const payload = {
                recipientIdType: 'meta_address',
                recipientK: recipientInfo.K,
                recipientV: recipientInfo.V,
                amount: 10,
            }

            const res = await axiosInstance.post('/send', payload);
            const { data } = res;

            console.log(res);
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }   
    })

    afterAll(async () => {
        try {
            await app.stop()
            await ganacheServer.close();
        } catch (err) {}
    })
});