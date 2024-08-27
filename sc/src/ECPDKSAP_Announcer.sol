// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

contract ECPDKSAP_Announcer {

  function sendEthViaContract(
    address payable _stealthAddress,
    bytes memory _R, 
    bytes memory _viewTag
  ) external payable {

    _announce(_stealthAddress, _R, _viewTag);

    _stealthAddress.transfer(msg.value);
  }

  function _announce(
    address _stealthAddress,
    bytes memory _R,
    bytes memory _viewTag
  ) internal {
    emit Announcement(SCHEME_ID, _stealthAddress, msg.sender, _R, _viewTag);
  }

  uint constant SCHEME_ID = 3327;

  event Announcement(
    uint256 indexed schemeId,
    address indexed stealthAddress,
    address indexed caller,
    bytes ephemeralPubKey,
    bytes metadata
  );
}
