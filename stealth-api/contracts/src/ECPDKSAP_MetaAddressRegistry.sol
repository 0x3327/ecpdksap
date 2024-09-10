// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import { IECPDKSAP_MetaAddressRegistry } from "./interface/IECPDKSAP_MetaAddressRegistry.sol";
import { ErrorCodes } from "./Utils.sol";

contract ECPDKSAP_MetaAddressRegistry is IECPDKSAP_MetaAddressRegistry {
  /// @inheritdoc IECPDKSAP_MetaAddressRegistry
  function registerMetaAddress(string memory _id, bytes memory _metaAddress) external payable {
    bytes32 _accessKey = keccak256(abi.encode(_id, "string"));

    require(s_idToMetaAddress[_accessKey].length == 0, ErrorCodes.META_ID_ALREADY_REGISTERED);

    s_idToMetaAddress[_accessKey] = _metaAddress;

    emit MetaAddressRegistered(_id, _metaAddress);
  }

  /// @inheritdoc IECPDKSAP_MetaAddressRegistry
  function resolve(string memory _id) external view returns (bytes memory metaAddress) {
    bytes32 _accessKey = keccak256(abi.encode(_id, "string"));

    metaAddress = s_idToMetaAddress[_accessKey];

    require(metaAddress.length != 0, ErrorCodes.META_ID_IS_NOT_REGISTERED);
  }

  /// @notice Maps the keccak256(_id) to the underlying meta address bytes
  mapping(bytes32 => bytes) public s_idToMetaAddress;
}
