// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

/// @notice Interface for calling the `ECPDKSAP_Announcer` contract which emits
///         information about eth transfers (Sender Ephermeral public key and view tag)
interface IECPDKSAP_Announcer {
  /// @notice Sends ether using the announcer contract as a proxy
  /// @param _stealthAddress Destination for the ether (Recipient's stealth address)
  /// @param _R Sender's epheremeral key
  /// @param _viewTag Protocol's view tag
  function sendEthViaProxy(address payable _stealthAddress, bytes memory _R, bytes memory _viewTag)
    external
    payable;

  /// @notice Notifies the recipient of the ether transfer without using the contract as a proxy
  /// @param _R Sender's epheremeral key
  /// @param _viewTag Protocol's view tag
  function ethSentWithoutProxy(bytes memory _R, bytes memory _viewTag) external;

  /// @notice Emitted when something is sent to a stealth address.
  /// @param schemeId Identifier corresponding to the applied stealth address scheme
  /// @param stealthAddress The computed stealth address for the recipient.
  /// @param caller The caller of the `announce` function that emitted this event.
  /// @param ephemeralPubKey Ephemeral public key used by the sender to derive the `stealthAddress`.
  /// @param metadata Arbitrary data to emit with the event. The first byte MUST be the view tag.
  event Announcement(
    uint256 indexed schemeId,
    address indexed stealthAddress,
    address indexed caller,
    bytes ephemeralPubKey,
    bytes metadata
  );
}
