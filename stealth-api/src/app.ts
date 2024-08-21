import { API } from '../services/api';
import BlockchainListener from '../services/blockchain-listener';

interface AppConfig {
    apiConfig: {
        serverName: string;
        host: string;
        port: string;
    };
}

class App {
    public config: AppConfig;
    public api!: API;
    public blockchainListener!: BlockchainListener;

    constructor(config: AppConfig) {
        this.config = config;
    }

    async start(): Promise<void> {
        // Load services
        // this.stealthService = new StealthService(this);
        this.api = new API(this);
        this.blockchainListener = new BlockchainListener(this);
        
        // Start API
        await this.api.start();

        const currentBLockNumber = await this.blockchainListener.getCurrentBlockNumber();
        const transactionDetails = await this.blockchainListener.getTransaction('0x16191fcc73ba807341cce0a93b32a63f66b5d28e1580d3ea402331951ced5bec');

        console.log('currentBlockNumber: ', currentBLockNumber);
        console.log('transactionDetails: ', transactionDetails);
    }
}

export default App;