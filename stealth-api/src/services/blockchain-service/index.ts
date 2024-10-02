import { Contract, ethers, Provider, Wallet } from 'ethers';
import App from '../../app';
import metaAddressArtifacts from '../../../artifacts/contracts/src/ECPDKSAP_MetaAddressRegistry.sol/ECPDKSAP_MetaAddressRegistry.json';
import announcerArtifacts from '../../../artifacts/contracts/src/ECPDKSAP_Announcer.sol/ECPDKSAP_Announcer.json';
import winston from 'winston';
import { groth16 } from 'snarkjs';
import fs from 'fs';

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
        this.listenVerifiedEvent();
        this.listenNullifierEvent();
        this.listenDebugProofEvent();
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

    public async registerMetaAddress(id: string, metaAddress: string) {
        try {
            const metaAddressBytes = ethers.toUtf8Bytes(metaAddress);

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
            return ethers.toUtf8String(metaAddress);
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
            const [K, V] = Buffer.from(R.slice(2, R.length), 'hex').toString('ascii').split('.');

            const amount = await this.provider.getBalance(stealthAddress);
            // console.log(event);

            // this.app.goHandler.receiveScan(

            // )

            this.logger.info('Announcement received:', schemaId, stealthAddress, sender, ethers.hexlify(R), ethers.hexlify(viewTag));
            await this.app.db.models.receivedTransactions.create({
                transaction_hash: event.log.transactionHash,
                block_number: event.log.blockNumber,
                amount,
                ephemeral_key: R,
                view_tag: viewTag,
                stealth_address: stealthAddress,
            });
            console.log('Announcement saved');
        });

        this.logger.info('Listening for Announcement event...');
    }

    public listenDebugProofEvent() {
        this.contracts.metaAddressRegistry.on("DebugProof", (_pA, _pB, _pC, _pubSignals, event) => {
            console.log("------------------------------");
            console.log("DebugProof event");
            console.log("pA: ", _pA);
            console.log("pB: ", _pB);
            console.log("pC: ", _pC);
            console.log("pubSignals: ", _pubSignals);
            console.log("Event: ", event);
        })
    }

    public listenVerifiedEvent() {
        this.contracts.metaAddressRegistry.on("ProofVerified", (result, event) => {
            console.log("--------------------------");
            console.log("Recieved ProofVerified event");
            console.log("Result: ", result);
            console.log("Event: ", event);
        });
    }

    public listenNullifierEvent() {
        this.contracts.metaAddressRegistry.on("NullifierRegistered", (nullifier, event) => {
            console.log("--------------------------");
            console.log("Recieved NullifierRegistered event");
            console.log("Nullifier: ", nullifier);
            console.log("Event: ", event);
        });
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

    public async verify(proof: any, publicSignals: any) {
        try {
            console.log("----------------------------");
            console.log("Data sending to contract");
            console.log(proof.pi_a);
            console.log(proof.pi_b);
            console.log(proof.pi_c);
            console.log(publicSignals);

            // const vKey = JSON.parse(fs.readFileSync("./verification_key.json").toString());
            // const result = await groth16.verify(vKey, publicSignals, proof);
            // console.log("Result: ", result);

            //const tx = await this.contracts.metaAddressRegistry.registerMetaAddress(
            //    "test",
            //    Buffer.from("testing", "utf-8"),
            //    proof.pi_a,
            //    proof.pi_b,
            //    proof.pi_c,
            //    publicSignals
            //);
            const tx = await this.contracts.metaAddressRegistry.registerMetaAddress(
                "test",
                "0x00",
                ["0x2ce9dbb039ca2f1b38fcc95bd7631da9b8c4fde817e58ed661534c7bf2c24e88", "0x1f9c665b466007b4fa0efbe867c2e1d92b9e7837a8f1adb9410fcbd6308b9d72"],
                [["0x0000e83316989c1eb35e034561f42046b7cb241b6455d9beade186a9380c8182", "0x28c923af57bdd02ddaafd797fbf87580648bcc8e999dc3a6dc0aaee6a9a12174"], ["0x170ae699f6ca720187f453f04d2783cfbf1efd466ec419eaa8bf1bfb8b4be81a", "0x058335ce5a3d511c86cd5dfb5265ef27a4995fe081bb3dc3db3617b16b4e19b3"]],
                ["0x101ab59c806005062464e7505715954fa0256620ae955443e04196fa6d0322da", "0x08ae2e25081369624765eac4f609768efb7e87937b1ab680122bf418c8670904"],
                ["0x12056cc40bb6e94c4264ba1ae08913ccc8ad0dfcaa6b7ec921be9c98305fac28", "0x0000000000000000000000000000000000000000000000000000000000002710", "0x25626f51857aecdb452fac9aec230ed0ee4c3cc077f9faa8a437079bd5a564ac"]
            );
            const recepit = await tx.wait();
            //console.log("Transaction: ", tx);
            //console.log("Receipt: ", recepit);
        }
        catch (error) {
            console.log("Error verifying proof: ", error);
        }
    }

    public async stop() {
        await this.contracts.metaAddressRegistry.off('MetaAddressRegistered');
        await this.contracts.announcer.off('Announcement');
        try {
            this.provider.destroy();
        } catch (err) {}
    }
}

export default BlockchainService;