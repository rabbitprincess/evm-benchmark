// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import {Script, console} from "forge-std/Script.sol";
import {Storage} from "../src/3_Storage.sol";

contract CounterScript is Script {
    Storage public s;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();

        s = new Storage();

        vm.stopBroadcast();
    }
}
