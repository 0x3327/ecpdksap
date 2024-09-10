// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Script.sol";

import { ECPDKSAP_Announcer } from "../src/ECPDKSAP_Announcer.sol";
import { ECPDKSAP_MetaAddressRegistry } from  "../src/ECPDKSAP_MetaAddressRegistry.sol";

contract DeploymentScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PK");
        vm.startBroadcast(deployerPrivateKey);

        ECPDKSAP_Announcer announcer = new ECPDKSAP_Announcer();

        ECPDKSAP_MetaAddressRegistry metaAddressRegistry = new ECPDKSAP_MetaAddressRegistry();

        vm.stopBroadcast();
    }
}

//forge script --chain sepolia script/NFT.s.sol:MyScript --rpc-url $SEPOLIA_RPC_URL --broadcast --verify -vvvv
