// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Hash{
    function hashSingle(uint256 x) external pure returns (bytes32) {
        return keccak256(abi.encodePacked(x));
    }

    function hashAssembly(uint256 x) external pure returns (bytes32 result) {
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, x)
            result := keccak256(ptr, 0x20)
        }
    }
}