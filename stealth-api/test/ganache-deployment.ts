import ganache, { Server } from "ganache";
import { ContractFactory, ethers } from 'ethers';

import anouncerArtifacts from "../artifacts/contracts/src/ECPDKSAP_Announcer.sol/ECPDKSAP_Announcer.json";
import metaAddressArtifacts from "../artifacts/contracts/src/ECPDKSAP_MetaAddressRegistry.sol/ECPDKSAP_MetaAddressRegistry.json";
import erc5564AnnouncerArtifacts from "../artifacts/contracts/src/ERC5564Announcer.sol/ERC5564Announcer.json";
import { Config } from "../types";
import configLoader from "../utils/config-loader";

type BlockchainParams = {
  ganacheServer: Server,
  privateKey: string;
  deployedContracts: {
    announcer: string;
    metaAddress: string;
  };
};

let config: Config = configLoader.load('test');

export async function deployContracts(): Promise<BlockchainParams> {
  return new Promise((resolve, reject) => {
    const mnemonic = 'test test test test test test test test test test test junk';
    const server = ganache.server({ wallet: { mnemonic }, logging: { quiet: true } });
    const port = 8545;

    server.listen(port, async (err) => {
        if (err) {
            reject(err);
        }

        console.log(`ganache listening on port ${server.address().port}...`);
        const provider = ethers.getDefaultProvider(`http://127.0.0.1:${port}`);

        const wallet = ethers.Wallet.fromPhrase(mnemonic);
        const account = wallet.connect(provider);
        config.stealthConfig.transferAddress = wallet.address;

        const erc5564AnnouncerFactory = new ContractFactory(erc5564AnnouncerArtifacts.abi, erc5564AnnouncerArtifacts.bytecode, account);
        const announcerFactory = new ContractFactory(anouncerArtifacts.abi, anouncerArtifacts.bytecode, account);
        const metaAddressFactory = new ContractFactory(metaAddressArtifacts.abi, metaAddressArtifacts.bytecode, account);

        console.log('Deploying meta address contract...');
        const metaAddress = (await (await metaAddressFactory.deploy({nonce: 0})).waitForDeployment());
        console.log('Deploying ERC5564 Announcer contract...');
        const erc5564Announcer = (await (await erc5564AnnouncerFactory.deploy({nonce: 1})).waitForDeployment());
        console.log('Deploying ECPDKSAP Announcer contract...');
        const announcer = (await (await announcerFactory.deploy(await erc5564Announcer.getAddress(), {nonce: 2})).waitForDeployment());

        console.log('Contracts deployed.')
        console.log('erc5564Announcer:', await erc5564Announcer.getAddress());
        console.log('metaAddress:', await metaAddress.getAddress());
        console.log('announcer:', await announcer.getAddress());
        

        resolve({
          ganacheServer: server,
            privateKey: '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80',
            deployedContracts: {
                announcer: await announcer.getAddress(),
                metaAddress: await metaAddress.getAddress()
            }
        })
    });
  });
}
