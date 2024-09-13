import { API } from './services/api';
import BlockchainService from './services/blockchain-service';
import GoHandler from './services/go-service';
import DB from './services/db';
import { Config } from '../types';
import { Op } from 'sequelize';
import { Info, ReceiveScanInfo, SendInfo } from './types';
import axios from 'axios';
import dotenv from 'dotenv';
import LoggerService from './services/logger';

dotenv.config({ path: `.env.development` });

require('../public/wasm_exec.js');

class App {
    public config: Config;
    public api!: API;
    public blockchainService!: BlockchainService;
    public goHandler!: GoHandler;
    public db!: DB;
    public loggerService!: LoggerService;

    constructor(config: Config) {
        this.config = config;

        // Initialize services
        this.loggerService = new LoggerService(this.config);
        this.goHandler = new GoHandler(this);
        this.api = new API(this);
        this.blockchainService = new BlockchainService(this);
        this.db = new DB(this.config.dbConfig);
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

        // const senderInfo: Info = await this.goHandler.genSenderInfo();
        // console.log("senderInfo", senderInfo);
        // const recipientInfo: Info = await this.goHandler.genRecipientInfo();
        // console.log("recipientInfo", recipientInfo);
        // const sendInfo: SendInfo = await this.goHandler.send(senderInfo.r, recipientInfo.K, recipientInfo.V);
        // console.log("sendInfo", sendInfo);
        // const receiveScanInfo: ReceiveScanInfo[] = await this.goHandler.receiveScan(recipientInfo.k, recipientInfo.v, [senderInfo.R], [sendInfo.viewTag]);
        // console.log("receiveScanInfo", receiveScanInfo)

        // // This will pass
        // this.db.models.receivedTransactions.create({
        //     transaction_hash: '0x123456',
        //     block_number: 3,
        //     amount: 100,
        //     stealth_address: '0xabcd12345678',
        //     ephemeral_key: '0x23456789098765432345678',
        //     view_tag: '0x23'
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
        //     recipient_k: '0x123450334565674',
        //     recipient_v: '0x123450abc431232',
        //     recipient_stealth_address: '0xabcdef54321',
        //     ephemeral_key: '0x12345789abc1232',
        // })

        // // Fetch all receipts by given parameters
        // const res_received = await this.db.models.receivedTransactions.findAll({
        //     where: {
        //         block_number: {
        //             [Op.between]: [1, 4],
        //         },
        //         amount: 100,
        //     },
        // });
        // // console.log(res_received);

        // const res_sent = await this.db.models.sentTransactions.findAll({
        //     where: {
        //         block_number: {
        //             [Op.between]: [1, 4],
        //         },
        //         amount: 101,
        //     },
        // });
        // console.log(res_sent);

        // Testirati funckcionalnosti BlockchainService-a
        // const receipt = await this.BlockchainService.registerMetaAddress('1', '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef');
        // console.log("app print", receipt);

        // const metaAddress = await this.BlockchainService.resolveMetaAddress('1');
        // console.log("app print", metaAddress);

        // const receipt = await this.BlockchainService.sendEthViaProxy(sendInfo.address, sendInfo.pubKey, sendInfo.viewTag, '0.001');
        // console.log("app print", receipt);
    }

    // async testRoutes(route: string): Promise<void> {
    //     console.log('Testing routes functionality using axios');
    //     console.log(`route: ${route}`);
    //     switch (route) {
    //         case '/':
    //             console.log('usao /');
    //             await axios.get(`http://localhost:${process.env.API_PORT || 8765}/`)
    //                 .then((res) => {
    //                     console.log('GET / response:', res.data);
    //                 })
    //                 .catch((err) => {
    //                     console.error('Error in GET /:', err.message);
    //                 });
    //             break;
    //         case '/send':
    //             console.log('usao /send');
    //             await axios.post(`http://localhost:${process.env.API_PORT || 8765}/send`, {
    //                 r: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', // it needs to be a valid r
    //                 K: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', // it needs to be a valid K
    //                 V: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', // it needs to be a valid V
    //             }).then((res) => {
    //                 console.log('POST /send response:', res.data);
    //             }).catch((err) => {
    //                 console.error('Error in POST /send:', err.message);
    //             });
    //             break;
    //         case '/check-received':
    //             console.log('usao /check-received');
    //             await axios.get(`http://localhost:${process.env.API_PORT || 8765}/check-received`, {
    //                 params: { // here maybe needs to be body and different params for input parameters for this route
    //                     fromBlock: 1, // it needs to be a valid fromBlock   
    //                     toBlock: 10, // it needs to be a valid toBlock 
    //                 },
    //             }).then((res) => {
    //                 console.log('GET /check-received response:', res.data);
    //             }).catch((err) => {
    //                 console.error('Error in GET /check-received:', err.message);
    //             });
    //             break;
    //         case '/transfer/:recieveId':
    //             console.log('usao /transfer/:recieveId');
    //             await axios.get(`http://localhost:${process.env.API_PORT || 8765}/transfer/:recieveId`, {
    //                 params: {
    //                     recieveId: 1, // it needs to be a valid recieveId
    //                     address: '0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef', // it needs to be a valid address    
    //                     amount: 100, // it needs to be a valid amount
    //                 },
    //             }).then((res) => {
    //                 console.log('GET /transfer response:', res.data);
    //             }).catch((err) => {
    //                 console.error('Error in GET /transfer:', err.message);
    //             });
    //             break;
    //     }
    // }
}

export default App;