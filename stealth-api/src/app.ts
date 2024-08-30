import { API } from './services/api';
import BlockchainListener from './services/blockchain-listener';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from './types';

require('../public/wasm_exec.js');

const { readFileSync } = require('fs');

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
        this.blockchainListener = new BlockchainListener(this);
        this.db = new DB(this.config.dbConfig);

        // Init database
        await this.db.sequelize.sync({ force: false });
        
        // Start API
        await this.api.start();

        const senderInfo: Info = await this.goHandler.genSenderInfo();
        console.log("senderInfo", senderInfo);
        const recipientInfo: Info = await this.goHandler.genRecipientInfo();
        console.log("recipientInfo", recipientInfo);
        const sendInfo: SendInfo = await this.goHandler.send(senderInfo.r, recipientInfo.K,recipientInfo.V);
        console.log("sendInfo", sendInfo);
        const receiveScanInfo: ReceiveScanInfo[] = await this.goHandler.receiveScan(senderInfo.k, senderInfo.v, [senderInfo.R], [sendInfo.viewTag]);
        console.log("receiveScanInfo", receiveScanInfo)

        // This will pass
        this.db.models.receivedTransactions.create({
            transaction_hash: '0x123456',
            block_number: 3,
            amount: 100,
            stealth_address: '0xabcd12345678',
            ephemeral_key: '0x23456789098765432345678',
            view_tag: '0x23'
        }).then(() => console.log('Created!'))

        // This will fetch the latest synced block
        this.db.models.receivedTransactions.max('block_number').then((res: number) => console.log(res))

        // This will create new sent transaction details, EC point coordinates separated by '.'
        this.db.models.sentTransactions.create({
            transaction_hash: '0x123456',
            block_number: 4,
            amount: 101,
            recipient_identifier: 'pera.eth',
            recipient_identifier_type: 'eth_dns',
            recipient_k: '0x12345.0x3345674',
            recipient_v: '0x12345.0xabc1232',
            recipient_stealth_address: '0xabcdef54321',
            ephemeral_key: '0x12345.0xabc1232',
        })

        // Fetch all receipts by given parameters
        const res_received = await this.db.models.receivedTransactions.findAll({
            block_number: {
                [Op.gte]: 1, // block_number >= 1
                [Op.lte]: 4, // block_number <= 4
                amount: 100  // amount == 100
            }
        });
        console.log(res_received);

        const res_sent = await this.db.models.sentTransactions.findAll({
            block_number: {
                [Op.gte]: 1, // block_number >= 1
                [Op.lte]: 4, // block_number <= 4
                amount: 101  // amount == 101
            }
        });
        console.log(res_sent);

        // Testirati funckcionalnosti BlockchainListener-a
    }
}

export default App;