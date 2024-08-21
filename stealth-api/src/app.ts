import { API } from '../services/api';
import BlockchainListener from '../services/blockchain-listener';
import GoHandler from '../services/go-service';

require('../public/wasm_exec.js');

import path from 'path';

const { readFileSync } = require('fs');

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
    public goHandler!: GoHandler;

    constructor(config: AppConfig) {
        this.config = config;
    }

    async start(): Promise<void> {
        // Load services
        // this.stealthService = new StealthService(this);
        this.goHandler = new GoHandler();
        this.api = new API(this);
        // this.blockchainListener = new BlockchainListener(this);

        const metaInfo = await this.goHandler.genSenderInfo();
        console.log(metaInfo);
        

        // Start API
        await this.api.start();

        // const currentBLockNumber = await this.blockchainListener.getCurrentBlockNumber();
        // const transactionDetails = await this.blockchainListener.getTransaction('0x16191fcc73ba807341cce0a93b32a63f66b5d28e1580d3ea402331951ced5bec');

        // console.log('currentBlockNumber: ', currentBLockNumber);
        // console.log('transactionDetails: ', transactionDetails);
    }
}

export default App;