import { describe, expect, test } from '@jest/globals';
import App from '../app';
import { after, before } from 'node:test';
import configLoader from '../../utils/config-loader';
import axios, { AxiosInstance } from 'axios';

// Application object
let app: App;

// Axios instance
let axiosInstance: AxiosInstance; 

describe('API routes test', () => {
    beforeAll(async () => {
        const config = configLoader.load('test');
        // console.log(config);

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

    afterAll(async () => {
        await app.stop();
    })
});