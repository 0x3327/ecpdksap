import dotenv from 'dotenv';
import { Config } from '../types';

const configLoader = {
    load(configType: string = 'development'): Config {
        dotenv.config({ path: `.env.${configType}` });

        const serverName = process.env.API_SERVER_NAME || 'defaultServerName';
        const host = process.env.API_HOST || '0.0.0.0';
        const port = process.env.API_PORT ?  parseInt(process.env.API_PORT) : 8765;

        const dbHost = process.env.DB_HOST || '0.0.0.0';
        const dbPort = process.env.DB_PORT ? parseInt(process.env.DB_PORT) : 3306;
        const dbDatabase = process.env.DB_DATABASE || 'stealthdb';
        const dbUsername = process.env.DB_USERNAME;
        const dbPassword = process.env.DB_PASSWORD;

        const config: Config = {
            apiConfig: {
                serverName,
                host,
                port,
            },
            dbConfig: {
                host: dbHost,
                port: dbPort,
                database: dbDatabase,
                username: dbUsername,
                password: dbPassword,
            }
        };

        return config;
    }
}

export default configLoader;