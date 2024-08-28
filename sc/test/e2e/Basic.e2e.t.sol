// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import { ECPDKSAP_Announcer } from "../../src/ECPDKSAP_Announcer.sol";
import { Constants } from "../../src/Utils.sol";

contract ECPDKSAP_Basic_E2E_Test is Test {
  function test_E2E_0() public {
    vm.startBroadcast(s_recipient);

    // TODO when the typescript API lib is ready
    // string[] memory inputs = new string[](3);
    // inputs[0] = "../impl/builds/ecpdksap-ll-latest-arm64";
    // inputs[1] = "bench";
    // inputs[2] = "only-bn254";
    // bytes memory res = vm.ffi(inputs);
    // string memory output = abi.decode(res, (string));

    // console.log(output);
    // vm.parseJson(output, "");

    vm.stopBroadcast();
  }

  function setUp() public {
    vm.deal(s_recipient, 100 ether);
  }

  address s_recipient = address(0x01ff);
}
