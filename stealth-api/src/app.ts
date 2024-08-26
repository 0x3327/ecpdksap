import { API } from './services/api';
import BlockchainListener from './services/blockchain-listener';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';

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

// type SendInfo = {
//     pubKey: string;
//     address: string;
//     viewTag: string;
// }

class App {
    public config: Config;
    public api!: API;
    public blockchainListener!: BlockchainListener;
    public goHandler!: GoHandler;
    public db!: DB;

    constructor(config: Config) {
        this.config = config;
    }

    async start(): Promise<void> {
        // Load services
        this.goHandler = new GoHandler();
        this.api = new API(this);
        // this.blockchainListener = new BlockchainListener(this);
        this.db = new DB(this.config.dbConfig);

        const senderInfo = await this.goHandler.genSenderInfo();
        console.log("senderInfo", senderInfo);
        const recipientInfo = await this.goHandler.genRecipientInfo();
        console.log("recipientInfo", recipientInfo);
        const sendInfo = await this.goHandler.send();
        console.log("sendInfo", sendInfo);
        const receiveScanInfo = await this.goHandler.receiveScan();
        console.log("receiveScanInfo", receiveScanInfo)

        // Init database
        await this.db.sequelize.sync({ force: false });
        
        // Start API
        await this.api.start();


        // TODO: ### UNCOMMENT, TEST, THEN REMOVE THIS !!! ###

        // // This will fail because of the missing parameters
        // this.db.models.receivedTransactions.create({
        //     transaction_hash: '0x123456',
        //     block_number: 3,
        //     amount: 100
        // }).then(() => console.log('Created!'))

        // // This will pass
        // this.db.models.receivedTransactions.create({
        //     transaction_hash: '0x123456',
        //     block_number: 3,
        //     amount: 100,
        //     stealthAddress: '0xabcd12345678',
        //     ephemeralKey: '0x23456789098765432345678',
        // }).then(() => console.log('Created!'))

        // // This will fetch the latest synced block
        // this.db.models.receivedTransactions.max('block_number').then((res: number) => console.log(res))

        // // This will create new sent transaction details, EC point coordinates separated by '.'
        // this.db.models.sentTransactions.create({
        //     transaction_hash: '0x123456',
        //     block_number: 4,
        //     amount: 101,
        //     recipient_identifier: 'pera.eth',
        //     recipient_identifier_type: 'eth_dns',
        //     recipient_k: '0x12345.0x3345674',
        //     recipient_v: '0x12345.0xabc1232',
        //     recipient_stealth_address: '0xabcdef54321',
        //     ephemeral_key: '0x12345.0xabc1232'
        // })



        // const currentBLockNumber = await this.blockchainListener.getCurrentBlockNumber();
        // const transactionDetails = await this.blockchainListener.getTransaction('0x16191fcc73ba807341cce0a93b32a63f66b5d28e1580d3ea402331951ced5bec');

        // console.log('currentBlockNumber: ', currentBLockNumber);
        // console.log('transactionDetails: ', transactionDetails);
    }
}

export default App;