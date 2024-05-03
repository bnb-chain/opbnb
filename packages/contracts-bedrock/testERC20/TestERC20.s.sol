// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "./TestERC20.sol";

contract DeployERC20 is Script {
    function run() external {
        vm.startBroadcast();
        TestERC20 erc20 = new TestERC20("My Token", "MTK", 100000000000000000000000);
        console.log("ERC20 Contract Address:", address(erc20));
        vm.stopBroadcast();
    }
}
