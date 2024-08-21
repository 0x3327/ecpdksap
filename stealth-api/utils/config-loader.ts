import dotenv from 'dotenv';

interface ApiConfig {
    serverName: string;
    host: string;
    port: string;
}

interface Config {
    apiConfig: ApiConfig;
}

const configLoader = {
    load(configType: string = 'development'): Config {
        dotenv.config({ path: `.env.${configType}` });

        const serverName = process.env.API_SERVER_NAME || 'defaultServerName';
        const host = process.env.API_HOST || '0.0.0.0';
        const port = process.env.API_PORT || '3000';

        const config: Config = {
            apiConfig: {
                serverName,
                host,
                port,
            },
        };

        return config;
    }
}

export default configLoader;