import { ethers } from 'ethers';
import dotenv from 'dotenv';
import App from '../../src/app';
import abi_registry from '../abi/MetaAddressRegistry.json';
import abi_announcer from '../abi/Announcer.json';

dotenv.config({ path: `.env.development` });

const provider = new ethers.InfuraProvider('sepolia', process.env.INFURA_API_KEY);
const privateKey = process.env.PRIVATE_KEY;
const wallet = new ethers.Wallet(privateKey!, provider);

const metaAddressRegistry = new ethers.Contract('0xB4B82918613524DB74967CA6c71979cD030B7991', abi_registry, wallet);
const announcer = new ethers.Contract('0x79820C9a124023D47BbCC6d0a24DB4D0075Ca724', abi_announcer, wallet);

class BlockchainListener {
    public app: App;

    constructor(app: App) {
        this.app = app;
        this.listenMetaAddressRegistredEvent();
        this.listenAnnouncementEvent();
    }

    public async getCurrentBlockNumber() {
        try {
            const blockNumber = await provider.getBlockNumber();
            return blockNumber;
        } catch (error) {
            console.error('Error fetching block number:', error);
        }
    }

    public async getTransaction(transactionHash: string) {
        try {
            const transaction = await provider.getTransaction(transactionHash);
            return transaction
        } catch (error) {
            console.error('Error fetching transaction details:', error);
        }
    }

    public async registerMetaAddress(id: string, metaAddress: string) {
        try {
            const metaAddressBytes = ethers.toUtf8Bytes(metaAddress);

            const tx = await metaAddressRegistry.registerMetaAddress(id, metaAddressBytes, {
                value: ethers.parseEther('0.00001')
            });
            console.log('Transaction sent:', tx.hash);

            const receipt = await tx.wait();
            console.log('Transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            console.error('Error registering meta address:', error);
        }
    }

    public async resolveMetaAddress(id: string) {
        try {
            const metaAddress = await metaAddressRegistry.resolve(id);
            console.log('Meta address resolved:', ethers.toUtf8String(metaAddress));
            return ethers.toUtf8String(metaAddress);
        } catch (error) {
            console.error('Error resolving meta address:', error);
        }
    }

    public async sendEthViaProxy(stealthAddress: string, R: string, viewTag: string, amount: string) {
        try {
            const tx = await announcer.sendEthViaProxy(stealthAddress, ethers.toUtf8Bytes(R), ethers.toUtf8Bytes(viewTag), {
                value: ethers.parseEther(amount)
            });
            console.log('Transaction sent:', tx.hash);

            const receipt = await tx.wait();
            console.log('Transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            console.error('Error sending ETH via proxy:', error);
        }
    }

    public async ethSentWithoutProxy(R: string, viewTag: string) {
        try {
            const tx = await announcer.ethSentWithoutProxy(ethers.toUtf8Bytes(R), ethers.toUtf8Bytes(viewTag));
            console.log('Transaction sent:', tx.hash);

            const receipt = await tx.wait();
            console.log('Transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            console.error('Error sending ETH without proxy:', error);
        }
    }

    public async listenMetaAddressRegistredEvent()  {
        metaAddressRegistry.on('MetaAddressRegistered', (id: string, metaAddress: string) => {
            console.log('Meta address registered:', id, ethers.hexlify(metaAddress));
        });

        console.log('Listening for MetaAddressRegistered event...');
    }

    public async listenAnnouncementEvent() {
        announcer.on('Announcement', (schemaId: string, stealthAddress: string, sender: string, R: string, viewTag: string) => {
            console.log('Announcement received:', schemaId, stealthAddress, sender, ethers.hexlify(R), ethers.hexlify(viewTag));
            this.app.db.models.sentTransactions.create({
                transaction_hash: '0x123456',
                block_number: 4,
                amount: 101,
                recipient_identifier: stealthAddress,
                recipient_identifier_type: null,
                recipient_k: '0x123450334565674',
                recipient_v: '0x123450abc431232',
                recipient_stealth_address: stealthAddress,
                ephemeral_key: R,
            });
        });

        console.log('Listening for Announcement event...');
    }

    public async transferEth(address: string, amount: string, privKey: string) {
        const signer = new ethers.Wallet(privKey, provider);
        try {
            const tx = await signer.sendTransaction({
                to: address,
                value: ethers.parseEther(amount)
            });
            console.log('Transaction sent:', tx.hash);

            const receipt = await tx.wait();
            return receipt;
        } catch (error) {
            console.error('Error sending ETH:', error);
        }
    }

    public async getBalance(address: string) {
        const balance = await provider.getBalance(address);
        return balance;
    }
}

export default BlockchainListener;