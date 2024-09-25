import { Contract, ethers, Provider, Wallet } from 'ethers';
import App from '../../app';
import metaAddressArtifacts from '../../../artifacts/contracts/src/ECPDKSAP_MetaAddressRegistry.sol/ECPDKSAP_MetaAddressRegistry.json';
import announcerArtifacts from '../../../artifacts/contracts/src/ECPDKSAP_Announcer.sol/ECPDKSAP_Announcer.json';
import winston, { add } from 'winston';

type Contracts = {
    metaAddressRegistry: Contract,
    announcer: Contract,
};

class BlockchainService {
    app: App;
    provider: Provider;
    wallet: Wallet;
    contracts: Contracts;
    logger: winston.Logger;

    constructor(app: App) {
        this.app = app;
        this.logger = app.loggerService.logger;

        // Load blockchain configuration
        const { privateKey, providerType, deployedContracts, infuraApiKey } = app.config.blockchainConfig;

        // Create provider
        switch (providerType) {
            case 'sepolia':
                this.provider = new ethers.InfuraProvider('sepolia', infuraApiKey);
                break;
            case 'ganache':
            default:
                this.provider = ethers.getDefaultProvider('http://localhost:8545');
        }

        // Create wallet
        this.wallet = new ethers.Wallet(privateKey, this.provider);

        // Connect contracts
        this.contracts = {
            metaAddressRegistry: new ethers.Contract(deployedContracts.metaAddress, metaAddressArtifacts.abi, this.wallet),
            announcer: new ethers.Contract(deployedContracts.announcer, announcerArtifacts.abi, this.wallet),
        }

        this.listenMetaAddressRegistredEvent();
        this.listenAnnouncementEvent();
    }

    public async getCurrentBlockNumber() {
        try {
            const blockNumber = await this.provider.getBlockNumber();
            return blockNumber;
        } catch (error) {
            this.logger.error('Error fetching block number:', error);
        }
    }

    public async getTransaction(transactionHash: string) {
        try {
            const transaction = await this.provider.getTransaction(transactionHash);
            return transaction;
        } catch (error) {
            this.logger.error('Error fetching transaction details:', error);
        }
    }

    public async registerMetaAddress(id: string, K: string, V: string) {
        try {
            const metaAddressBytes = `0x${Buffer.from(K + '|' + V, 'ascii').toString('hex')}`

            const tx = await this.contracts.metaAddressRegistry.registerMetaAddress(id, metaAddressBytes, {
                value: ethers.parseEther('0.00001')
            });

            this.logger.info('Meta address registration, transaction sent:', tx.hash);
            const receipt = await tx.wait();

            this.logger.info('Meta address registration, transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            this.logger.error('Error registering meta address:', error);
        }
    }

    public async resolveMetaAddress(id: string) {
        try {
            const metaAddress = await this.contracts.metaAddressRegistry.resolve(id);
            this.logger.info('Meta address resolved:', ethers.toUtf8String(metaAddress));

            const metaAddressString = ethers.toUtf8String(metaAddress)
            const rawMetaAddress = metaAddressString.split('|');
            return {
                K: rawMetaAddress[0],
                V: rawMetaAddress[1],
            }
        } catch (error) {
            this.logger.error('Error resolving meta address:', error);
        }
    }

    public async sendEthViaProxy(stealthAddress: string, R: string, viewTag: string, amount: string) {
        try {
            console.log({
                stealthAddress, 
                R: `0x${Buffer.from(R, 'ascii').toString('hex')}`,
                viewTag: `0x${viewTag}`
            })
            const tx = await this.contracts.announcer.sendEthViaProxy(stealthAddress, `0x${Buffer.from(R, 'ascii').toString('hex')}`, `0x${viewTag}`, {
                value: ethers.parseEther(amount)
            });
            this.logger.info('Sending ETH via Proxy, transaction sent:', tx.hash);

            const receipt = await tx.wait();
            this.logger.info('Sending ETH via Proxy, transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            console.log(error);
            this.logger.error('Error sending ETH via proxy:', error);
        }
    }

    public async ethSentWithoutProxy(stealthAddress: string, R: string, viewTag: string, amount: string) {
        try {
            const tx = await this.contracts.announcer.ethSentWithoutProxy(ethers.toUtf8Bytes(R), ethers.toUtf8Bytes(viewTag));
            await this.wallet.sendTransaction({
                to: stealthAddress,
                value: ethers.parseEther(amount),
              });
              
            this.logger.info('Sending ETH Without Proxy, Transaction sent:', tx.hash);

            const receipt = await tx.wait();
            this.logger.info('Sending ETH Without Proxy, transaction confirmed:', receipt.transactionHash);
            return receipt;
        } catch (error) {
            this.logger.error('Error sending ETH without proxy:', error);
        }
    }

    public async listenMetaAddressRegistredEvent()  {
        this.contracts.metaAddressRegistry.on('MetaAddressRegistered', (id: string, metaAddress: string) => {
            this.logger.info('Meta address registered:', id, ethers.hexlify(metaAddress));
        });

        this.logger.info('Listening for MetaAddressRegistered event...');
    }

    public async listenAnnouncementEvent() {
        this.contracts.announcer.on('Announcement', async (...parameters) => {
            const [schemaId, stealthAddress, sender, R, viewTag, event] = parameters;
            // console.log('Announcement received:', schemaId, stealthAddress, sender, ethers.hexlify(R), ethers.hexlify(viewTag))
            const cleanedR = R.slice(2);
            const RGo = Buffer.from(cleanedR, "hex").toString("ascii");
            const viewTagGo = viewTag.slice(2);

            const amount = await this.provider.getBalance(stealthAddress);
            // console.log(event);
            const matched = await this.app.goHandler.receiveScan(
                this.app.config.stealthConfig.k,
                this.app.config.stealthConfig.v,
                [RGo],
                [viewTagGo]
            );
            
            if (matched.length === 0) {
                return;
            }

            this.logger.info('Announcement received:', schemaId, stealthAddress, sender, ethers.hexlify(R), ethers.hexlify(viewTag));
            const computedAddress = matched[0].address;

            if (computedAddress === stealthAddress.toLowerCase()) {   
                await this.app.db.models.receivedTransactions.create({
                    transaction_hash: event.log.transactionHash,
                    block_number: event.log.blockNumber,
                    amount,
                    ephemeral_key: RGo,
                    view_tag: viewTagGo,
                    stealth_address: computedAddress,
                });
                console.log('Announcement saved');
            }
        });

        this.logger.info('Listening for Announcement event...');
    }

    public async transferEth(address: string, amount: string, privKey: string) {
        const signer = new ethers.Wallet(privKey, this.provider);
        try {
            const tx = await signer.sendTransaction({
                to: address,
                value: ethers.parseEther(amount)
            });
            this.logger.info('Transfer ETH, transaction sent:', tx.hash);

            const receipt = await tx.wait();
            return receipt;
        } catch (error) {
            this.logger.error('Error sending ETH:', error);
        }
    }

    public async getBalance(address: string) {
        const balance = await this.provider.getBalance(address);
        return balance;
    }

    public async stop() {
        await this.contracts.metaAddressRegistry.off('MetaAddressRegistered');
        await this.contracts.announcer.off('Announcement');

        // Add a delay to ensure all pending operations are completed
        await new Promise(resolve => setTimeout(resolve, 1000));

        try {
            this.provider.destroy();
        } catch (err) {}
    }
}

export default BlockchainService;