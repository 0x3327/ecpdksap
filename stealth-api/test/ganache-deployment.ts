import ganache, { Server } from "ganache";
import { ContractFactory, ethers } from 'ethers';

import anouncerArtifacts from "../artifacts/contracts/ECPDKSAP_Announcer.sol/ECPDKSAP_Announcer.json";
import metaAddressArtifacts from "../artifacts/contracts/ECPDKSAP_MetaAddressRegistry.sol/ECPDKSAP_MetaAddressRegistry.json";

type BlockchainParams = {
  ganacheServer: Server,
  privateKey: string;
  deployedContracts: {
    announcer: string;
    metaAddress: string;
  };
};

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
        const provider = new ethers.JsonRpcProvider(`http://127.0.0.1:${port}`);

        const wallet = ethers.Wallet.fromPhrase(mnemonic);
        const account = wallet.connect(provider);

        const announcerFactory = new ContractFactory(anouncerArtifacts.abi, anouncerArtifacts.bytecode, account);
        const metaAddressFactory = new ContractFactory(metaAddressArtifacts.abi, metaAddressArtifacts.bytecode, account);

        const announcer = await announcerFactory.deploy();
        const metaAddress = await metaAddressFactory.deploy();

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
