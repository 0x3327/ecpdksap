// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import { ECPDKSAP_Announcer } from "../../src/ECPDKSAP_Announcer.sol";
import { IECPDKSAP_Announcer } from "../../src/interface/IECPDKSAP_Announcer.sol";

import { Constants } from "../../src/Utils.sol";

contract ECPDKSAP_Announcer_Test is Test {
  function testFuzz_sendEthViaProxy_ShouldEmitCorrectEvent(
    address payable _stealthAddress,
    bytes memory _R,
    bytes memory _viewTag
  ) public {
    vm.startBroadcast(s_sender);

    vm.assume(_viewTag.length == 1 || _viewTag.length == 2);

    vm.expectEmit(true, true, true, false, address(s_announcer));
    emit IECPDKSAP_Announcer.Announcement(
      Constants.ECPDKSAP_SCHEME_ID, _stealthAddress, s_sender, _R, _viewTag
    );

    s_announcer.sendEthViaProxy(_stealthAddress, _R, _viewTag);

    vm.stopBroadcast();
  }

  function testFuzz_sendEthViaProxy_ShouldTransferEthCorrectly(
    address payable _stealthAddress,
    bytes memory _R,
    bytes memory _viewTag,
    uint256 _value
  ) public {
    vm.startBroadcast(s_sender);

    vm.assume(_value < 90 ether);

    vm.assume(_viewTag.length == 1 || _viewTag.length == 2);

    uint256 _stealthBalanceBefore = _stealthAddress.balance;

    s_announcer.sendEthViaProxy{ value: _value }(_stealthAddress, _R, _viewTag);

    uint256 _stealthBalanceAfter = _stealthAddress.balance;

    assertEq(_value, _stealthBalanceAfter - _stealthBalanceBefore);

    //note: there shouldn't be any residual funds at announcer contract
    assertEq(address(s_announcer).balance, 0);

    vm.stopBroadcast();
  }

  function testFuzz_ethSentWithoutProxy_ShouldEmitCorrectEvent(
    bytes memory _R,
    bytes memory _viewTag
  ) public {
    vm.startBroadcast(s_sender);

    vm.assume(_viewTag.length == 1 || _viewTag.length == 2);

    vm.expectEmit(true, false, true, false, address(s_announcer));
    emit IECPDKSAP_Announcer.Announcement(
      Constants.ECPDKSAP_SCHEME_ID, address(0x0), s_sender, _R, _viewTag
    );

    s_announcer.ethSentWithoutProxy(_R, _viewTag);

    vm.stopBroadcast();
  }

  //note: this test shouldn't be implemented, since there shouldn't be a Eth trail connection with the announcer contract
  // function testFuzz_ethSentWithoutProxy_ShouldTransferEthCorrectyl(bytes memory _R, bytes memory _viewTag, uint _value) public {}

  function setUp() public {
    s_announcer = new ECPDKSAP_Announcer();

    vm.deal(s_sender, 100 ether);
  }

  ECPDKSAP_Announcer s_announcer;

  address s_sender = address(0x00ff);
}
