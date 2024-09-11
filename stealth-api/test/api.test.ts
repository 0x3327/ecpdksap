import { describe, expect, test } from '@jest/globals';
import App from '../src/app';
import configLoader from '../utils/config-loader';
import axios, { AxiosInstance } from 'axios';
import { deployContracts } from './ganache-deployment';
import { Server } from 'ganache';
import { Config } from '../types';

// Application object
let app: App;

// Axios instance
let axiosInstance: AxiosInstance;

let ganacheServer: Server

let config: Config;

describe('API routes test', () => {
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

        // Start application
        await app.start();

        // Read API config params
        const { host, port } = app.config.apiConfig;

        // Initialize axios instance
        axiosInstance = axios.create({
            baseURL: `http://${host}:${port}`
          });

        // console.log(axiosInstance);
    }, 30000)

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

    test('Send stealth transaction via Proxy', async () => {
        console.log("----------------------- SEND ------------------------");
        try {
            const recipientInfo = await app.goHandler.genRecipientInfo();

            config.stealthConfig.recipientk = recipientInfo.k;
            config.stealthConfig.recipientv = recipientInfo.v;
            config.stealthConfig.recipientK = recipientInfo.K;
            config.stealthConfig.recipientV = recipientInfo.V;
            
            const payload = {
                recipientIdType: 'meta_address',
                recipientK: recipientInfo.K,
                recipientV: recipientInfo.V,
                amount: 10,
                withProxy: true,
            }

            const res = await axiosInstance.post('/send', payload);

            // TODO: Check response
            // console.log(res);

            // Wait for Announcement event
            await (new Promise((resolve, reject) => setTimeout(resolve, 5000)));

        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }   
    }, 30000);

    test('Check received funds', async () => {
        console.log("----------------------- CHECK-RECEIVED ------------------------");
        try {
            const res = await axiosInstance.get('/check-received');
            // console.log("res check-received", res);

            // Wait for Announcement event ???
            // await (new Promise((resolve, reject) => setTimeout(resolve, 20000)));
        } catch(err) {
            console.log(err);
            expect(true).toBe(false);
        }
    }, 30000);

    test('Transfer funds', async () => {
        console.log("----------------------- TRANSFER ------------------------");
        try {
            const res = await axiosInstance.get('/transfer/1');
            console.log("res transfer", res);
            // await (new Promise((resolve, reject) => setTimeout(resolve, 20000))); ???
        } catch (err) {
            console.log(err);
            expect(true).toBe(false);
        }
    }, 30000);

    afterAll(async () => {
        try {
            await app.stop()
            await ganacheServer.close();
        } catch (err) {}
    })
});
