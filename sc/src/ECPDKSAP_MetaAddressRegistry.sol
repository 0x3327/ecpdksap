// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

contract ECPDKSAP_MetaAddressRegistry {

  function register(
    string memory _id,
    bytes memory _spendingPubKey,
    bytes memory _viewingPubKey
  ) external payable {

    s_spendingPubKeys[msg.sender].push(_spendingPubKey);
    s_viewingPubKeys[msg.sender].push(_viewingPubKey);

    require(s_idToPubKeys[_id].length == 0, "Meta addr. `id` already used!");

    s_idToPubKeys[_id].push([_spendingPubKey, _viewingPubKey]);

    emit MetaAddressRegistered(_id, _spendingPubKey, _viewingPubKey);
  }

    mapping(address => bytes[]) public s_spendingPubKeys;
    mapping(address => bytes[]) public s_viewingPubKeys;

    mapping(string => bytes[][]) public s_idToPubKeys;

    event MetaAddressRegistered (string id, bytes spendingPubKey, bytes viewingPubKey);
}
