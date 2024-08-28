// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import { IERC5564Announcer } from "./interface/IERC5564Announcer.sol";
import { IECPDKSAP_Announcer } from "./interface/IECPDKSAP_Announcer.sol";
import { Constants } from "./Utils.sol";

contract ECPDKSAP_Announcer is IECPDKSAP_Announcer {
  /// @inheritdoc IECPDKSAP_Announcer
  function sendEthViaProxy(address payable _stealthAddress, bytes memory _R, bytes memory _viewTag)
    external
    payable
  {
    _announce(_stealthAddress, _R, _viewTag);

    _stealthAddress.transfer(msg.value);
  }

  /// @inheritdoc IECPDKSAP_Announcer
  function ethSentWithoutProxy(bytes memory _R, bytes memory _viewTag) external {
    _announce(address(0x0), _R, _viewTag);
  }

  /// @notice Called by integrators to emit an `Announcement` event and interact with the Singleton contract
  /// @param _stealthAddress The computed stealth address for the recipient.
  /// @param _R Sender's ephermeral public key
  /// @param _viewTag Protocol's view tag
  function _announce(address _stealthAddress, bytes memory _R, bytes memory _viewTag) internal {
    emit Announcement(Constants.ECPDKSAP_SCHEME_ID, _stealthAddress, msg.sender, _R, _viewTag);

    IERC5564Announcer(Constants.ERC5564_ANNOUNCER_ADDRESS).announce(
      Constants.ECPDKSAP_SCHEME_ID, _stealthAddress, _R, _viewTag
    );
  }
}
