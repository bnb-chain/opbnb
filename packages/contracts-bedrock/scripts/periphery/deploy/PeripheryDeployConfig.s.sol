// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { Script } from "forge-std/Script.sol";
import { console2 as console } from "forge-std/console2.sol";
import { stdJson } from "forge-std/StdJson.sol";

/// @title PeripheryDeployConfig
/// @notice Represents the configuration required to deploy the periphery contracts. It is expected
///         to read the file from JSON. A future improvement would be to have fallback
///         values if they are not defined in the JSON themselves.
contract PeripheryDeployConfig is Script {
    string internal _json;

    // General configuration.
    string public create2DeploymentSalt;

    // Configuration for Gelato.
    address public gelatoAutomateContract;

    // Configuration for standard operations Drippie contract.
    address public operationsDrippieOwner;

    // Configuration for the faucet Drippie contract.
    address public faucetDrippieOwner;

    // Configuration for the Faucet contract.
    address public faucetAdmin;
    address public faucetOnchainAuthModuleAdmin;
    uint256 public faucetOnchainAuthModuleTtl;
    uint256 public faucetOnchainAuthModuleAmount;
    address public faucetOffchainAuthModuleAdmin;
    uint256 public faucetOffchainAuthModuleTtl;
    uint256 public faucetOffchainAuthModuleAmount;
    bool public installOpChainFaucetsDrips;
    bool public archivePreviousOpChainFaucetsDrips;
    uint256 public smallOpChainFaucetDripValue;
    uint256 public smallOpChainFaucetDripInterval;
    uint256 public largeOpChainFaucetDripValue;
    uint256 public largeOpChainFaucetDripInterval;
    uint256 public opChainAdminWalletDripValue;
    uint256 public opChainAdminWalletDripInterval;
    address public opL1BridgeAddress;
    address public baseL1BridgeAddress;
    address public zoraL1BridgeAddress;
    address public pgnL1BridgeAddress;
    address public orderlyL1BridgeAddress;
    address public modeL1BridgeAddress;
    address public lyraL1BridgeAddress;
    address public liskL1BridgeAddress;
    address[6] public smallFaucetsL1BridgeAddresses;
    address[2] public largeFaucetsL1BridgeAddresses;
    uint256 public dripVersion;
    uint256 public previousDripVersion;

    // Configuration booleans.
    bool public deployDripchecks;
    bool public deployFaucetContracts;
    bool public deployOperationsContracts;

    constructor(string memory _path) {
        console.log("PeripheryDeployConfig: reading file %s", _path);
        try vm.readFile(_path) returns (string memory data) {
            _json = data;
        } catch {
            console.log("Warning: unable to read config. Do not deploy unless you are not using config.");
            return;
        }

        // General configuration.
        create2DeploymentSalt = stdJson.readString(_json, "$.create2DeploymentSalt");

        // Configuration for Gelato.
        gelatoAutomateContract = stdJson.readAddress(_json, "$.gelatoAutomateContract");

        // Configuration for the standard operations Drippie contract.
        operationsDrippieOwner = stdJson.readAddress(_json, "$.operationsDrippieOwner");

        // Configuration for the faucet Drippie contract.
        faucetDrippieOwner = stdJson.readAddress(_json, "$.faucetDrippieOwner");

        // Configuration for the Faucet contract.
        faucetAdmin = stdJson.readAddress(_json, "$.faucetAdmin");
        faucetOnchainAuthModuleAdmin = stdJson.readAddress(_json, "$.faucetOnchainAuthModuleAdmin");
        faucetOnchainAuthModuleTtl = stdJson.readUint(_json, "$.faucetOnchainAuthModuleTtl");
        faucetOnchainAuthModuleAmount = stdJson.readUint(_json, "$.faucetOnchainAuthModuleAmount");
        faucetOffchainAuthModuleAdmin = stdJson.readAddress(_json, "$.faucetOffchainAuthModuleAdmin");
        faucetOffchainAuthModuleTtl = stdJson.readUint(_json, "$.faucetOffchainAuthModuleTtl");
        faucetOffchainAuthModuleAmount = stdJson.readUint(_json, "$.faucetOffchainAuthModuleAmount");

        // Configuration booleans.
        deployDripchecks = stdJson.readBool(_json, "$.deployDripchecks");
        deployFaucetContracts = stdJson.readBool(_json, "$.deployFaucetContracts");
        deployOperationsContracts = stdJson.readBool(_json, "$.deployOperationsContracts");
    }
}
