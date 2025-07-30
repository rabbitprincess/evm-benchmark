// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Arithmetic {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        return a + b;
    }

    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        return a - b;
    }

    function mul(uint256 a, uint256 b) external pure returns (uint256) {
        return a * b;
    }

    function div(uint256 a, uint256 b) external pure returns (uint256) {
        return a / b;
    }
}