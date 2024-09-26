import { API } from './services/api';
import BlockchainService from './services/blockchain-service';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import dotenv from 'dotenv';
import LoggerService from './services/logger';
import CLIService from './services/cli';
import SocketService from './services/sockets';

dotenv.config({ path: `.env.development` });

require('../public/wasm_exec.js');

class App {
    public config: Config;
    public api!: API;
    public blockchainService!: BlockchainService;
    public goHandler!: GoHandler;
    public db!: DB;
    public loggerService!: LoggerService;
    public cliService!: CLIService;
    public socketService!: SocketService;

    constructor(config: Config) {
        this.config = config;

        // Initialize services
        this.loggerService = new LoggerService(this.config);
        this.goHandler = new GoHandler(this);
        this.api = new API(this);
        this.blockchainService = new BlockchainService(this);
        this.db = new DB(this.config.dbConfig);
        this.cliService = new CLIService(this);
        this.socketService = new SocketService(this);
    }

    async stop(): Promise<void> {
        this.loggerService.logger.info('Service stopping');
        await this.api.stop();
        await this.blockchainService.stop();
        this.loggerService.logger.info('App stopped');
        this.socketService.stop();
    }

    async start(): Promise<void> {

        // Init database
        await this.db.sequelize.sync({ force: true }); // true if you want to drop the table
        
        // Start API
        await this.api.start();

        // Start Sockets API
        this.socketService.start();
    }
}

export default App;