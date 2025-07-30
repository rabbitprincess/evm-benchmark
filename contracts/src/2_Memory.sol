// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Memory {
    function memorySetGetFree(uint256 x) external pure returns (uint256) {
        assembly {
            let ptr := mload(0x40)        // free ptr
            mstore(ptr, x)                // set
            let y := mload(ptr)           // get
            mstore(0x40, add(ptr, 0x20))  // free

            // return value
            mstore(0x00, y)
            return(0x00, 0x20)
        }
    }

    function memorySetGetFreeLoop(uint256 n) external pure returns (uint256) {
        assembly {
            let ptr := mload(0x40)    // free ptr
            let y := 0

            for { let i := 0 } lt(i, n) { i := add(i, 1) } {
                mstore(ptr, i)        // set
                y := mload(ptr)       // get
                ptr := add(ptr, 0x20) // move free pointer
            }
            mstore(0x40, ptr)

            // return value
            mstore(0x00, y)
            return(0x00, 0x20)
        }
    }
}
