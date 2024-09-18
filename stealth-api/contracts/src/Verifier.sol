// SPDX-License-Identifier: CC0-1.0
pragma solidity ^0.8.20;

contract NullifierRegistry {
    // Mapping to track registered nullifiers
    mapping(uint256 => bool) private nullifiers;

    // Event to emit when a nullifier is successfully registered
    event NullifierRegistered(uint256 nullifier);

    /**
     * @dev Registers a new nullifier if it does not exist.
     * Reverts if the nullifier already exists.
     * @param nullifier The nullifier to be registered.
     */
    function registerNullifier(uint256 nullifier) external {
        require(!nullifiers[nullifier], "Error: Nullifier already exists!");

        // Register the nullifier by setting its value to true
        nullifiers[nullifier] = true;

        // Emit an event after successful registration
        emit NullifierRegistered(nullifier);
    }

    /**
     * @dev Checks if a nullifier has been registered.
     * @param nullifier The nullifier to check.
     * @return bool indicating if the nullifier exists.
     */
    function isNullifierRegistered(uint256 nullifier) external view returns (bool) {
        return nullifiers[nullifier];
    }
}
