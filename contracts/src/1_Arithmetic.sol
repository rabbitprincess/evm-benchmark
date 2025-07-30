// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Arithmetic {
    function add(uint256 a, uint256 b, uint256 loop) external pure returns (uint256) {
        uint256 result = 0;
        for (uint256 i = 0; i < loop; i++) {
            result += a + b;
        }
        return result;
    }

    function sub(uint256 a, uint256 b, uint256 loop) external pure returns (uint256) {
        uint256 result = 0;
        for (uint256 i = 0; i < loop; i++) {
            result = a - b;
        }
        return result;
    }

    function mul(uint256 a, uint256 b, uint256 loop) external pure returns (uint256) {
        uint256 result = 0;
        for (uint256 i = 0; i < loop; i++) {
            result = a * b;
        }
        return result;
    }

    function div(uint256 a, uint256 b, uint256 loop) external pure returns (uint256) {
        uint256 result = 0;
        for (uint256 i = 0; i < loop; i++) {
            result = a / b;
        }
        return result;
    }
}