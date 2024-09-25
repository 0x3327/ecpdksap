import { API } from './services/api';
import BlockchainService from './services/blockchain-service';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import dotenv from 'dotenv';
import LoggerService from './services/logger';
import { CommandHandler } from './services/cli';
import SocketHandler from './services/sockets';
import { Server as SocketServer } from 'socket.io'
import { createServer, Server } from 'http';

dotenv.config({ path: `.env.development` });

require('../public/wasm_exec.js');

class App {
    public config: Config;
    public api!: API;
    public blockchainService!: BlockchainService;
    public goHandler!: GoHandler;
    public db!: DB;
    public loggerService!: LoggerService;
    public commandHandler!: CommandHandler;
    public httpServer!: Server;
    public socketServer!: SocketServer;
    public socketHandler!: SocketHandler;

    constructor(config: Config) {
        this.config = config;

        // Initialize services
        this.loggerService = new LoggerService(this.config);
        this.goHandler = new GoHandler(this);
        this.api = new API(this);
        this.blockchainService = new BlockchainService(this);
        this.db = new DB(this.config.dbConfig);
        this.commandHandler = new CommandHandler(this);
        this.httpServer = createServer();
        this.socketServer = new SocketServer(this.httpServer);
        this.socketHandler = new SocketHandler(this.socketServer, this);
    }

    async stop(): Promise<void> {
        await this.api.stop();
        await this.blockchainService.stop();
        console.log('App stopped')
    }

    async start(): Promise<void> {

        // Init database
        await this.db.sequelize.sync({ force: true }); // true if you want to drop the table
        
        // Start API
        await this.api.start();

        // Start Sockets API
    }
}

export default App;