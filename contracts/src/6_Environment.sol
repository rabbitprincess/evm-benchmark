// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract EnvironmentOps {
    /// @notice Return basic block and tx environment info
    function envBasic() external view returns (
        uint256 blockNumber,
        uint256 timestamp,
        uint256 chainId,
        address caller
    ) {
        blockNumber = block.number;
        timestamp = block.timestamp;
        chainId = block.chainid;
        caller = msg.sender;
    }

    /// @notice Return external account info (costlier)
    function envExternal(address target) external view returns (
        uint256 balance,
        uint256 codeSize
    ) {
        balance = target.balance;
        assembly {
            codeSize := extcodesize(target)
        }
    }

    /// @notice Return call value
    function envCallValue() external payable returns (uint256) {
        return msg.value;
    }
}