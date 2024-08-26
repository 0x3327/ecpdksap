export type Config = {
    apiConfig: ApiConfig;
    dbConfig: DbConfig;
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