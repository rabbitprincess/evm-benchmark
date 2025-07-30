// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract JumpTable {
    /// @notice Switch-based jump table benchmark
    function jumpSwitch(uint256 selector) external pure returns (uint256) {
        assembly {
            switch selector
            case 0 { mstore(0x00, 111) }
            case 1 { mstore(0x00, 222) }
            case 2 { mstore(0x00, 333) }
            case 3 { mstore(0x00, 444) }
            default { mstore(0x00, 999) }
            return(0x00, 0x20)
        }
    }

    /// @notice Loop with conditional jump to measure repeated JUMPI cost
    function jumpLoop(uint256 n) external pure returns (uint256) {
        assembly {
            let sum := 0
            for { let i := 0 } lt(i, n) { i := add(i, 1) } {
                if iszero(mod(i, 2)) {
                    sum := add(sum, 1)
                }
            }
            mstore(0x00, sum)
            return(0x00, 0x20)
        }
    }
}