// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import {Test, console} from "forge-std/Test.sol";
import {Storage} from "../src/3_Storage.sol";

contract CounterTest is Test {
    Storage public s;

    function setUp() public {
        s = new Storage();
        s.storageSetGet(0); 
    }

    function test_Increment() public {
        s.increment(1);
        assertEq(s.getValue(), 1);
    }

}
