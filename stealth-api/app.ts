import { API } from './services/api';
import { StealthService } from './services/stealth-service';
import { ethers } from 'ethers';

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
    public stealthService!: StealthService;

    constructor(config: AppConfig) {
        this.config = config;
    }

    async start(): Promise<void> {
        // Load services
        this.stealthService = new StealthService(this);
        this.api = new API(this);
        
        // Start API
        await this.api.start();
    }
}

export default App;