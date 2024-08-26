import { API } from '../services/api';
import BlockchainListener from '../services/blockchain-listener';
import GoHandler from '../services/go-service';

require('../public/wasm_exec.js');

const { readFileSync } = require('fs');

type Info = {
    k: string;
    v: string;
    r: string;
    K: string;
    V: string;
    R: string;
    P_sender: string;
    viewTag: string;
    P_Recipient: string;
    Version: string;
    ViewTagVersion: string;
};

type SendInfo = {
    pubKey: string;
    address: string;
    viewTag: string;
}

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
        this.goHandler = new GoHandler();
        this.api = new API(this);
        this.blockchainListener = new BlockchainListener(this);

        const senderInfo = await this.goHandler.genSenderInfo();
        console.log("senderInfo", senderInfo);
        const recipientInfo = await this.goHandler.genRecipientInfo();
        console.log("recipientInfo", recipientInfo);
        const sendInfo = await this.goHandler.send();
        console.log("sendInfo", sendInfo);
        const receiveScanInfo = await this.goHandler.receiveScan();
        console.log("receiveScanInfo", receiveScanInfo)
        
        // Start API
        await this.api.start();

        const currentBLockNumber = await this.blockchainListener.getCurrentBlockNumber();
        const transactionDetails = await this.blockchainListener.getTransaction('0x16191fcc73ba807341cce0a93b32a63f66b5d28e1580d3ea402331951ced5bec');

        // console.log('currentBlockNumber: ', currentBLockNumber);
        // console.log('transactionDetails: ', transactionDetails);
    }
}

export default App;