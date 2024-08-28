// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

/// @notice Contains all constants used in the project
library Constants {
  uint256 constant ECPDKSAP_SCHEME_ID = 3327;
  address constant ERC5564_ANNOUNCER_ADDRESS = address(0x55649E01B5Df198D18D95b5cc5051630cfD45564);
}

/// @notice Contains all possible error codes that could cause reverts
library ErrorCodes {
  string constant META_ID_ALREADY_REGISTERED = "ERR: Meta ID is already registered!";
  string constant META_ID_IS_NOT_REGISTERED = "ERR: Meta ID has not been registered!";
}
