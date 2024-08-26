import { ethers } from 'ethers';
import dotenv from 'dotenv';
import App from '../../src/app';

dotenv.config({ path: `.env.development` });

const provider = new ethers.WebSocketProvider(`wss://mainnet.infura.io/ws/v3/${process.env.INFURA_API_KEY}`);

class BlockchainListener {
    public app: App;

    constructor(app: App) {
        this.app = app;
    }

    public async getCurrentBlockNumber() {
        try {
            const blockNumber = await provider.getBlockNumber();
            // console.log('Current block number:', blockNumber);
            return blockNumber;
        } catch (error) {
            console.error('Error fetching block number:', error);
        }
    }

    public async getTransaction(transactionHash: string) {
        try {
            const transaction = await provider.getTransaction(transactionHash);
            // console.log('Transaction details:', transaction);
            return transaction
        } catch (error) {
            console.error('Error fetching transaction details:', error);
        }
    }
}

export default BlockchainListener;