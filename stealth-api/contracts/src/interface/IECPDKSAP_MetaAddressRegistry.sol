// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

/// @notice Interface for calling the `ECPDKSAP_MetaAddressRegistry` contract, which stores
/// information about used `id`s and corresponding Meta addresses (Spending & Viewing public keys)
interface IECPDKSAP_MetaAddressRegistry {
  /// @notice Registers an `_id` to the underlying meta address
  /// @param _id Identifier corresponding to the raw bytes `_metaAddress`
  /// @param _metaAddress Encoded Spending and Viewing public keys
  /// @param _nullifier Nullifier for check
  /// @dev function is `payable` to the possibility of introducing registration fees in the future
  function registerMetaAddress(string memory _id, bytes memory _metaAddress, uint256 _nullifier) external payable;

  /// @notice Resolves an `_id` to the underlying meta address
  /// @param _id Identifier corresponding to the raw bytes `_metaAddress`
  /// @return metaAddress Corresponding encoded Spending and Viewing public keys
  function resolve(string memory _id) external view returns (bytes memory metaAddress);

  /// @notice Emitted when a Meta address is registered
  /// @param id Identifier corresponding to the raw bytes `_metaAddress`
  /// @param metaAddress Encoded Spending and Viewing public keys
  event MetaAddressRegistered(string indexed id, bytes indexed metaAddress);
}
