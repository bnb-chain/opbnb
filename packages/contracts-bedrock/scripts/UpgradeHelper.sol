// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { VmSafe } from "forge-std/Vm.sol";
import { Script } from "forge-std/Script.sol";

import { console2 as console } from "forge-std/console2.sol";
import { stdJson } from "forge-std/StdJson.sol";

import { Deployer } from "scripts/Deployer.sol";

import { Proxy } from "src/universal/Proxy.sol";
import { L1StandardBridge } from "src/L1/L1StandardBridge.sol";
import { StandardBridge } from "src/universal/StandardBridge.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { L1CrossDomainMessenger } from "src/L1/L1CrossDomainMessenger.sol";
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { OptimismMintableERC20Factory } from "src/universal/OptimismMintableERC20Factory.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { Constants } from "src/libraries/Constants.sol";
import { L1ERC721Bridge } from "src/L1/L1ERC721Bridge.sol";
import { ProtocolVersions, ProtocolVersion } from "src/L1/ProtocolVersions.sol";
import { StorageSetter } from "src/universal/StorageSetter.sol";
import { ProxyAdmin } from "src/universal/ProxyAdmin.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Chains } from "scripts/Chains.sol";
import { Config } from "scripts/Config.sol";

import { ChainAssertions } from "scripts/ChainAssertions.sol";
import { Types } from "scripts/Types.sol";
import { LibStateDiff } from "scripts/libraries/LibStateDiff.sol";
import { EIP1967Helper } from "test/mocks/EIP1967Helper.sol";
import { ForgeArtifacts } from "scripts/ForgeArtifacts.sol";
import { Process } from "scripts/libraries/Process.sol";

/// @title Upgrade
/// @notice Upgrade used to help upgrade opBNB contracts.
contract UpgradeHelper {

    /// @notice Semantic version.
    /// @custom:semver 1.0.0
    string public constant version = "1.0.1";

    function upgrade(address payable _proxyAdmin, address payable[] calldata _proxys, address[] calldata _implementations) public {
        require(_proxys.length == _implementations.length, "proxy length not equal impl");
        bytes memory _innerCallData = abi.encodeCall(SuperchainConfig.initialize, (_proxyAdmin, false));
        (bool success, bytes memory data) = _proxyAdmin.delegatecall(
            abi.encodeWithSignature("upgradeAndCall(address,address,bytes)", _proxys[0], _implementations[0], _innerCallData)
        );
        require(success, "superchainconfig upgrade failed");
    }

}
