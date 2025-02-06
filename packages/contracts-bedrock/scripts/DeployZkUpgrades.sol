// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { VmSafe } from "forge-std/Vm.sol";
import { Script } from "forge-std/Script.sol";

import { console2 as console } from "forge-std/console2.sol";
import { stdJson } from "forge-std/StdJson.sol";

import { Deployer } from "scripts/Deployer.sol";

import { Proxy } from "src/universal/Proxy.sol";
import { OptimismPortal2 } from "src/L1/OptimismPortal2.sol";
import { Constants } from "src/libraries/Constants.sol";
import { Chains } from "scripts/Chains.sol";
import { Config } from "scripts/Config.sol";

import { ChainAssertions } from "scripts/ChainAssertions.sol";
import { Types } from "scripts/Types.sol";
import { LibStateDiff } from "scripts/libraries/LibStateDiff.sol";
import { EIP1967Helper } from "test/mocks/EIP1967Helper.sol";
import { ForgeArtifacts } from "scripts/ForgeArtifacts.sol";
import { Process } from "scripts/libraries/Process.sol";
import { AnchorStateRegistry } from "../src/dispute/AnchorStateRegistry.sol";
import { ZkFaultDisputeGame } from "../src/dispute/ZkFaultDisputeGame.sol";
import { SP1VerifierGateway } from "sp1-contracts/contracts/src/SP1VerifierGateway.sol";

/// @title Deploy
contract Deploy is Deployer {
    using stdJson for string;

    ////////////////////////////////////////////////////////////////
    //                        Modifiers                           //
    ////////////////////////////////////////////////////////////////

    /// @notice Modifier that wraps a function in broadcasting.
    modifier broadcast() {
        vm.startBroadcast(msg.sender);
        _;
        vm.stopBroadcast();
    }

    ////////////////////////////////////////////////////////////////
    //                        Accessors                           //
    ////////////////////////////////////////////////////////////////

    /// @notice The create2 salt used for deployment of the contract implementations.
    ///         Using this helps to reduce config across networks as the implementation
    ///         addresses will be the same across networks when deployed with create2.
    function _implSalt() internal view returns (bytes32) {
        return keccak256(bytes(Config.implSalt()));
    }

    /// @notice Returns the proxy addresses, not reverting if any are unset.
    function _proxiesUnstrict() internal view returns (Types.ContractSet memory proxies_) {
        proxies_ = Types.ContractSet({
            DisputeGameFactory: getAddress("DisputeGameFactoryProxy"),
            DelayedWETH: getAddress("DelayedWETHProxy"),
            AnchorStateRegistry: getAddress("AnchorStateRegistryProxy"),
            ZkFaultProofConfig: getAddress("ZkFaultProofConfigProxy")
        });
    }

    ////////////////////////////////////////////////////////////////
    //                    SetUp and Run                           //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy all of the zk contracts.
    function run() public {
        console.log("Deploy all of the zk contracts");
        _run();
    }

    /// @notice Internal function containing the deploy logic.
    function _run() internal virtual {
        console.log("start of deploying zk contracts!");
        deployProxies();
        deployImplementations();
        initializeImplementations();
    }

    /// @notice Deploy all of the proxies
    function deployProxies() public {
        console.log("Deploying proxies");
        deployERC1967Proxy("DisputeGameFactoryProxy");
        deployERC1967Proxy("DelayedWETHProxy");
        deployERC1967Proxy("AnchorStateRegistryProxy");
        deployERC1967Proxy("ZkFaultProofConfigProxy");
    }

    /// @notice Deploy all of the implementations
    function deployImplementations() public {
        console.log("Deploying implementations");
        deployDisputeGameFactory();
        deployDelayedWETH();
        deployAnchorStateRegistry();
        deployZkFaultProofConfig();
        deploySp1VerifierSuite();
        deployZkFaultDisputeGame();
        deployOptimismPortal2();
    }

    /// @notice Initialize all of the implementations
    function initializeImplementations() public {
        console.log("Initializing implementations");

        initializeDisputeGameFactory();
        initializeDelayedWETH();
        initializeAnchorStateRegistry();
        initializeZkFaultProofConfig();
    }

    ////////////////////////////////////////////////////////////////
    //                Proxy Deployment Functions                  //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploys an ERC1967Proxy contract with the ProxyAdmin as the owner.
    /// @param _name The name of the proxy contract to be deployed.
    /// @return addr_ The address of the deployed proxy contract.
    function deployERC1967Proxy(string memory _name) public returns (address addr_) {
        address proxyAdmin = msg.sender;
        console.log("%s proxyAdmin at %s", _name, address(proxyAdmin));
        addr_ = deployERC1967ProxyWithOwner(_name, proxyAdmin);
    }

    /// @notice Deploys an ERC1967Proxy contract with a specified owner.
    /// @param _name The name of the proxy contract to be deployed.
    /// @param _proxyOwner The address of the owner of the proxy contract.
    /// @return addr_ The address of the deployed proxy contract.
    function deployERC1967ProxyWithOwner(
        string memory _name,
        address _proxyOwner
    )
    public
    broadcast
    returns (address addr_)
    {
        console.log(string.concat("Deploying ERC1967 proxy for ", _name));
        Proxy proxy = new Proxy({ _admin: _proxyOwner });

        require(EIP1967Helper.getAdmin(address(proxy)) == _proxyOwner);

        save(_name, address(proxy));
        console.log("   at %s", address(proxy));
        addr_ = address(proxy);
    }

    ////////////////////////////////////////////////////////////////
    //             Implementation Deployment Functions            //
    ////////////////////////////////////////////////////////////////

    /// @notice Deploy the DisputeGameFactory
    function deployDisputeGameFactory() public broadcast returns (address addr_) {
        console.log("Deploying DisputeGameFactory implementation");
        DisputeGameFactory factory = new DisputeGameFactory{ salt: _implSalt() }();
        save("DisputeGameFactory", address(factory));
        console.log("DisputeGameFactory deployed at %s", address(factory));

        addr_ = address(factory);
    }

    function deployDelayedWETH() public broadcast returns (address addr_) {
        console.log("Deploying DelayedWETH implementation");
        DelayedWETH weth = new DelayedWETH{ salt: _implSalt() }(cfg.faultGameWithdrawalDelay());
        save("DelayedWETH", address(weth));
        console.log("DelayedWETH deployed at %s", address(weth));

        addr_ = address(weth);
    }

    /// @notice Deploy the AnchorStateRegistry
    function deployAnchorStateRegistry() public broadcast returns (address addr_) {
        console.log("Deploying AnchorStateRegistry implementation");
        AnchorStateRegistry anchorStateRegistry =
                    new AnchorStateRegistry{ salt: _implSalt() }(DisputeGameFactory(mustGetAddress("DisputeGameFactoryProxy")));
        save("AnchorStateRegistry", address(anchorStateRegistry));
        console.log("AnchorStateRegistry deployed at %s", address(anchorStateRegistry));

        addr_ = address(anchorStateRegistry);
    }

    function deployZkFaultProofConfig() public broadcast returns (address addr_) {
        console.log("Deploying ZkFaultProofConfig implementation");
        ZkFaultProofConfig config = new ZkFaultProofConfig{ salt: _implSalt() }();
        save("ZkFaultProofConfig", address(config));
        console.log("ZkFaultProofConfig deployed at %s", address(config));

        addr_ = address(config);
    }

    function deploySp1VerifierSuite() public broadcast returns (address addr_) {
        console.log("Deploying SP1VerifierGateway implementation");
        SP1VerifierGateway suite = new SP1VerifierGateway{ salt: _implSalt() }();
        save("SP1VerifierGateway", address(suite));
        console.log("SP1VerifierGateway deployed at %s", address(suite));


        addr_ = address(suite);
        console.log("Deploying PlonkVerifier implementation");

        address verifier = address(new SP1Verifier());
        suite.addRoute(verifier);
    }

    function deployZkFaultDisputeGame() public broadcast returns (address addr_) {
        console.log("Deploying ZkFaultDisputeGame implementation");
        ZkFaultDisputeGame zkFaultDisputeGame = new ZkFaultDisputeGame{ salt: _implSalt() }({
            _gameType: GameType.wrap(uint32(3)),
            _maxGenerateProofDuration: Duration.wrap(uint64(86400)),
            _maxDetectFaultDuration: Duration.wrap(uint64(86400)),
            _PROPOSER_BOND: 1 ether,
            _CHALLENGER_BOND: 1 ether,
            _FEE_VAULT_ADDRESS: address(0),
            _CHALLENGER_REWARD_PERCENTAGE: 1000,
            _PROVER_REWARD_PERCENTAGE: 5000,
            _weth: IDelayedWETH(0x4aa45Ad189d7dd1Cf595f93867d02e5DD3Fb603c),
            _anchorStateRegistry: IAnchorStateRegistry(0x2f11AE10A3DBaCFE33ce5d058CF71C68feB73f40),
            _config: ZkFaultProofConfig(0xAa75d486f7fD9BE4bCe16111946dFf6eF6303536),
            _l2ChainId: 901
        });
        save("ZkFaultDisputeGame", address(zkFaultDisputeGame));
        console.log("ZkFaultDisputeGame deployed at %s", address(zkFaultDisputeGame));

        addr_ = address(zkFaultDisputeGame);
    }

    /// @notice Deploy the OptimismPortal2
    function deployOptimismPortal2() public broadcast returns (address addr_) {
        console.log("Deploying OptimismPortal2 implementation");

        OptimismPortal2 portal = new OptimismPortal2{ salt: _implSalt() }({
            _proofMaturityDelaySeconds: cfg.proofMaturityDelaySeconds(),
            _disputeGameFinalityDelaySeconds: cfg.disputeGameFinalityDelaySeconds()
        });

        save("OptimismPortal2", address(portal));
        console.log("OptimismPortal2 deployed at %s", address(portal));

        addr_ = address(portal);
    }

    /// @notice Initialize the DisputeGameFactory
    function initializeDisputeGameFactory() public broadcast {
        console.log("Upgrading and initializing DisputeGameFactory proxy");
        address disputeGameFactoryProxy = mustGetAddress("DisputeGameFactoryProxy");
        address disputeGameFactory = mustGetAddress("DisputeGameFactory");

        Proxy proxy = Proxy(mustGetAddress("DisputeGameFactoryProxy"));
        string data = abi.encodeCall(DisputeGameFactory.initialize, Constants.BSCQANET_PROXY_ADMIN);
        Proxy.upgradeToAndCall({
            _implementation: disputeGameFactory,
            _data: data
        });

        string memory version = DisputeGameFactory(disputeGameFactoryProxy).version();
        console.log("DisputeGameFactory version: %s", version);

        ChainAssertions.checkDisputeGameFactory({ _contracts: _proxiesUnstrict(), _expectedOwner: Constants.BSCQANET_PROXY_ADMIN });
    }

    function initializeDelayedWETH() public broadcast {
        console.log("Upgrading and initializing DelayedWETH proxy");
        address delayedWETHProxy = mustGetAddress("DelayedWETHProxy");
        address delayedWETH = mustGetAddress("DelayedWETH");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        _upgradeAndCallViaSafe({
            _proxy: payable(delayedWETHProxy),
            _implementation: delayedWETH,
            _innerCallData: abi.encodeCall(DelayedWETH.initialize, (msg.sender, SuperchainConfig(superchainConfigProxy)))
        });

        string memory version = DelayedWETH(payable(delayedWETHProxy)).version();
        console.log("DelayedWETH version: %s", version);

        ChainAssertions.checkDelayedWETH({
            _contracts: _proxiesUnstrict(),
            _cfg: cfg,
            _isProxy: true,
            _expectedOwner: msg.sender
        });
    }

}
