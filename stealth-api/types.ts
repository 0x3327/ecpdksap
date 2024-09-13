export type Config = {
    logging: boolean,
    apiConfig: ApiConfig;
    dbConfig: DbConfig;
    blockchainConfig: BlockchainConfig;
    stealthConfig: StealthConfig;
}

export type ApiConfig = {
    serverName: string;
    host: string;
    port: number;
}

export type DbConfig = {
    host: string;
    port: number;
    database: string;
    username: string | undefined;
    password: string | undefined
}

export type BlockchainConfig = {
    privateKey: string,
    providerType: string, 
    deployedContracts: {
        announcer: string,
        metaAddress: string,
    }, 
    infuraApiKey?: string
}

export type StealthConfig = {
    k: string,
    v: string,
}