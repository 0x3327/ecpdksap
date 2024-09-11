import dotenv from 'dotenv';
import { Config } from '../types';

const configLoader = {
    load(configType: string = 'test'): Config {
        dotenv.config({ path: `.env.${configType}` });

        const serverName = process.env.API_SERVER_NAME || 'defaultServerName';
        const host = process.env.API_HOST || '0.0.0.0';
        const port = process.env.API_PORT ?  parseInt(process.env.API_PORT) : 8765;

        const dbHost = process.env.DB_HOST || '0.0.0.0';
        const dbPort = process.env.DB_PORT ? parseInt(process.env.DB_PORT) : 3306;
        const dbDatabase = process.env.DB_DATABASE || 'stealthdb';
        const dbUsername = process.env.DB_USERNAME;
        const dbPassword = process.env.DB_PASSWORD;

        const privateKey = process.env.BLOCKCHAIN_PRIVATE_KEY!;
        const providerType = process.env.BLOCKCHAIN_PROVIDER_TYPE!;
        const announcer = process.env.BLOCKCHAIN_CONTRACT_ANNOUNCER!;
        const metaAddress = process.env.BLOCKCHAIN_CONTRACT_META_ADDRESS!;

        const logging = process.env.LOGGING === 'true';
        
        const senderRandomness = process.env.SENDER_r!;
        const senderR = process.env.SENDER_R!;
        const recipientk = process.env.RECIPIENT_k!;
        const recipientK = process.env.RECIPIENT_K!;
        const recipientv = process.env.RECIPIENT_v!;
        const recipientV = process.env.RECIPIENT_V!;
        const Rs = process.env.RS ? process.env.RS.split(';') : [];
        const ViewTags = process.env.VIEW_TAG ? process.env.VIEW_TAG.split(';') : [];
        const transferAddress = process.env.TRANSFER_ADDRESS!;

        const config: Config = {
            logging,
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
            },
            blockchainConfig: {
                privateKey,
                providerType,
                deployedContracts: {
                    announcer,
                    metaAddress,
                }
            },
            stealthConfig:  {
                senderRandomness,
                senderR,
                recipientk,
                recipientK,
                recipientv,
                recipientV,
                Rs,
                ViewTags,
                transferAddress
            }
        };

        return config;
    }
}

export default configLoader;