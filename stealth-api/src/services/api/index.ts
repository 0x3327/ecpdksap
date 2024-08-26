import express, { Express, Request, Response } from 'express';
import bodyParser from 'body-parser';
import routeHandlers from './route-handlers';
import App from '../../app';
import { queryParser } from 'express-query-parser';

interface AppConfig {
    serverName: string;
    host: string;
    port: string;
}

class API {
    private app: App;
    private host: string;
    private port : number;
    private serverName: string;
    private server: Express;

    constructor(app: App) {
        const { serverName, host, port } = app.config.apiConfig;
        this.app = app;
        this.host = host;
        this.port = port;
        this.serverName = serverName;

        this.server = express();
        this.server.use(bodyParser.json());
        this.server.use(queryParser({
            parseNull: true,
            parseUndefined: true,
            parseBoolean: true,
            parseNumber: true
          }));

        this.exposeRoutes();
    }
    
    private exposeRoutes(): void {
        const routes = routeHandlers(this.app);

        for (const routeHandler of routes) {
            const { method, handler, path } = routeHandler;

            console.log(`Exposing route: ${method.toUpperCase()}:${path}`)

            switch (method.toLowerCase()) {
                case 'post':
                    this.server.post(path, handler);
                    break;
                case 'get':
                    this.server.get(path, handler);
                    break;
                default:
                    throw new Error(`Unsupported method: ${method}`)
            }
        }
    }

    public start(): Promise<void> {
        return new Promise((resolve, reject) => {
            this.server.listen(this.port, this.host, () => {
                console.log(`${this.serverName} server started on ${this.host}:${this.port}`);
                resolve();
            });
        });
    }
}

export { API };