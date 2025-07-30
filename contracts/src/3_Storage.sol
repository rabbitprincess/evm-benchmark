// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

contract Storage {
    uint256 private value;

    /// @notice Single SSTORE followed by a single SLOAD
    function storageSetGet(uint256 x) external returns (uint256) {
        value = x;          // SSTORE
        return value;       // SLOAD
    }

    /// @notice Repeated SSTORE + SLOAD on the same slot
    function storageSetGetLoop(uint256 n) external returns (uint256) {
        uint256 last;
        for (uint256 i = 0; i < n; i++) {
            value = i;      // SSTORE
            last = value;   // SLOAD
        }
        return last;
    }

    /// @notice Pure SLOAD loop to measure storage caching effect
    /// First SLOAD costs ~2100 gas, subsequent cached reads cost ~100 gas
    function storageLoadLoop(uint256 n) external view returns (uint256) {
        uint256 last;
        for (uint256 i = 0; i < n; i++) {
            last = value;   // SLOAD
        }
        return last;
    }

    function increment(uint256 amount) external {
        value = value + amount;
    }

    function decrement(uint256 amount) external {
        value = value - amount;
    }
}
