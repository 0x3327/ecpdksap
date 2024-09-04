import { API } from './services/api';
import BlockchainListener from './services/blockchain-listener';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from './types';

require('../public/wasm_exec.js');

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
        const sendInfo: SendInfo = await this.goHandler.send(senderInfo.r, recipientInfo.K, recipientInfo.V);
        console.log("sendInfo", sendInfo);
        const receiveScanInfo: ReceiveScanInfo[] = await this.goHandler.receiveScan(recipientInfo.k, recipientInfo.v, [senderInfo.R], [sendInfo.viewTag]);
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
            recipient_k: '0x123450334565674',
            recipient_v: '0x123450abc431232',
            recipient_stealth_address: '0xabcdef54321',
            ephemeral_key: '0x12345789abc1232',
        })

        // Fetch all receipts by given parameters
        const res_received = await this.db.models.receivedTransactions.findAll({
            where: {
                block_number: {
                    [Op.between]: [1, 4],
                },
                amount: 100,
            },
        });
        console.log(res_received);

        const res_sent = await this.db.models.sentTransactions.findAll({
            where: {
                block_number: {
                    [Op.between]: [1, 4],
                },
                amount: 101,
            },
        });
        console.log(res_sent);

        // Testirati funckcionalnosti BlockchainListener-a
        // const receipt = await this.blockchainListener.registerMetaAddress('1', '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef');
        // console.log("app print", receipt);

        // const metaAddress = await this.blockchainListener.resolveMetaAddress('1');
        // console.log("app print", metaAddress);

        const receipt = await this.blockchainListener.sendEthViaProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, '0.001');
        console.log("app print", receipt);
        console.log("app print to", receipt.to);
        console.log("app print from", receipt.from);
        console.log("app print blockNumber", receipt.blockNumber);
        console.log("app print hash", receipt.hash);
    }
}

export default App;