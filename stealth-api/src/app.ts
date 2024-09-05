import { API } from './services/api';
import BlockchainListener from './services/blockchain-listener';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from './types';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config({ path: `.env.development` });

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
        // console.log(res_received);

        const res_sent = await this.db.models.sentTransactions.findAll({
            where: {
                block_number: {
                    [Op.between]: [1, 4],
                },
                amount: 101,
            },
        });
        // console.log(res_sent);

        // Testirati funckcionalnosti BlockchainListener-a
        // const receipt = await this.blockchainListener.registerMetaAddress('1', '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef');
        // console.log("app print", receipt);

        // const metaAddress = await this.blockchainListener.resolveMetaAddress('1');
        // console.log("app print", metaAddress);

        // const receipt = await this.blockchainListener.sendEthViaProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, '0.001');
        // console.log("app print", receipt);
    }

    async testRoutes(route: string) {
        console.log('Testing routes functionality using axios');
        console.log(`route: ${route}`);
        switch (route) {
            case '/':
                await axios.get(`http://localhost:${process.env.API_PORT || 8765}/`)
                    .then((res) => {
                        console.log('GET / response:', res.data);
                    })
                    .catch((err) => {
                        console.error('Error in GET /:', err.message);
                    });
            case '/send':
                await axios.post(`http://localhost:${process.env.API_PORT || 8765}/send`, {
                    r: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
                    K: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
                    V: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
                }).then((res) => {
                    console.log('POST /send response:', res.data);
                }).catch((err) => {
                    console.error('Error in POST /send:', err.message);
                });
            case '/check-received':
                await axios.get(`http://localhost:${process.env.API_PORT || 8765}/check-received`, {
                    params: {
                        fromBlock: 1,
                        toBlock: 10,
                    },
                }).then((res) => {
                    console.log('GET /check-received response:', res.data);
                }).catch((err) => {
                    console.error('Error in GET /check-received:', err.message);
                });
            case '/transfer/:recieveId':
                await axios.get(`http://localhost:${process.env.API_PORT || 8765}/transfer/:recieveId`, {
                    params: {
                        recieveId: 1,
                        address: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef',
                        amount: 100,
                    },
                }).then((res) => {
                    console.log('GET /transfer response:', res.data);
                }).catch((err) => {
                    console.error('Error in GET /transfer:', err.message);
                });
        }
    }
}

export default App;