// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/* Testing utilities */
import { Test, StdUtils } from "forge-std/Test.sol";
import { L2OutputOracle } from "../L1/L2OutputOracle.sol";
import { L2ToL1MessagePasser } from "../L2/L2ToL1MessagePasser.sol";
import { L1StandardBridge } from "../L1/L1StandardBridge.sol";
import { L2StandardBridge } from "../L2/L2StandardBridge.sol";
import { L1ERC721Bridge } from "../L1/L1ERC721Bridge.sol";
import { L2ERC721Bridge } from "../L2/L2ERC721Bridge.sol";
import { OptimismMintableERC20Factory } from "../universal/OptimismMintableERC20Factory.sol";
import { OptimismMintableERC721Factory } from "../universal/OptimismMintableERC721Factory.sol";
import { OptimismMintableERC20 } from "../universal/OptimismMintableERC20.sol";
import { OptimismPortal } from "../L1/OptimismPortal.sol";
import { L1CrossDomainMessenger } from "../L1/L1CrossDomainMessenger.sol";
import { L2CrossDomainMessenger } from "../L2/L2CrossDomainMessenger.sol";
import { SequencerFeeVault } from "../L2/SequencerFeeVault.sol";
import { FeeVault } from "../universal/FeeVault.sol";
import { AddressAliasHelper } from "../vendor/AddressAliasHelper.sol";
import { LegacyERC20ETH } from "../legacy/LegacyERC20ETH.sol";
import { Predeploys } from "../libraries/Predeploys.sol";
import { Types } from "../libraries/Types.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import { Proxy } from "../universal/Proxy.sol";
import { Initializable } from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import { ResolvedDelegateProxy } from "../legacy/ResolvedDelegateProxy.sol";
import { AddressManager } from "../legacy/AddressManager.sol";
import { L1ChugSplashProxy } from "../legacy/L1ChugSplashProxy.sol";
import { IL1ChugSplashDeployer } from "../legacy/L1ChugSplashProxy.sol";
import { Strings } from "@openzeppelin/contracts/utils/Strings.sol";
import { LegacyMintableERC20 } from "../legacy/LegacyMintableERC20.sol";
import { SystemConfig } from "../L1/SystemConfig.sol";
import { ResourceMetering } from "../L1/ResourceMetering.sol";
import { Constants } from "../libraries/Constants.sol";

contract CommonTest is Test {
    address alice = address(128);
    address bob = address(256);
    address multisig = address(512);

    address immutable ZERO_ADDRESS = address(0);
    address immutable NON_ZERO_ADDRESS = address(1);
    uint256 immutable NON_ZERO_VALUE = 100;
    uint256 immutable ZERO_VALUE = 0;
    uint64 immutable NON_ZERO_GASLIMIT = 50000;
    bytes32 nonZeroHash = keccak256(abi.encode("NON_ZERO"));
    bytes NON_ZERO_DATA = hex"0000111122223333444455556666777788889999aaaabbbbccccddddeeeeffff0000";

    event TransactionDeposited(
        address indexed from,
        address indexed to,
        uint256 indexed version,
        bytes opaqueData
    );

    FFIInterface ffi;

    function setUp() public virtual {
        // Give alice and bob some ETH
        vm.deal(alice, 1 << 16);
        vm.deal(bob, 1 << 16);
        vm.deal(multisig, 1 << 16);

        vm.label(alice, "alice");
        vm.label(bob, "bob");
        vm.label(multisig, "multisig");

        // Make sure we have a non-zero base fee
        vm.fee(1000000000);

        ffi = new FFIInterface();
    }

    function emitTransactionDeposited(
        address _from,
        address _to,
        uint256 _mint,
        uint256 _value,
        uint64 _gasLimit,
        bool _isCreation,
        bytes memory _data
    ) internal {
        emit TransactionDeposited(
            _from,
            _to,
            0,
            abi.encodePacked(_mint, _value, _gasLimit, _isCreation, _data)
        );
    }
}

contract L2OutputOracle_Initializer is CommonTest {
    // Test target
    L2OutputOracle oracle;
    L2OutputOracle oracleImpl;

    L2ToL1MessagePasser messagePasser =
        L2ToL1MessagePasser(payable(Predeploys.L2_TO_L1_MESSAGE_PASSER));

    // Constructor arguments
    address internal proposer = 0x000000000000000000000000000000000000AbBa;
    address internal owner = 0x000000000000000000000000000000000000ACDC;
    uint256 internal submissionInterval = 1800;
    uint256 internal l2BlockTime = 2;
    uint256 internal startingBlockNumber = 200;
    uint256 internal startingTimestamp = 1000;
    address guardian;

    // Test data
    uint256 initL1Time;

    event OutputProposed(
        bytes32 indexed outputRoot,
        uint256 indexed l2OutputIndex,
        uint256 indexed l2BlockNumber,
        uint256 l1Timestamp
    );

    event OutputsDeleted(uint256 indexed prevNextOutputIndex, uint256 indexed newNextOutputIndex);

    // Advance the evm's time to meet the L2OutputOracle's requirements for proposeL2Output
    function warpToProposeTime(uint256 _nextBlockNumber) public {
        vm.warp(oracle.computeL2Timestamp(_nextBlockNumber) + 1);
    }

    /// @dev Helper function to propose an output.
    function proposeAnotherOutput() public {
        bytes32 proposedOutput2 = keccak256(abi.encode());
        uint256 nextBlockNumber = oracle.nextBlockNumber();
        uint256 nextOutputIndex = oracle.nextOutputIndex();
        warpToProposeTime(nextBlockNumber);
        uint256 proposedNumber = oracle.latestBlockNumber();

        // Ensure the submissionInterval is enforced
        assertEq(nextBlockNumber, proposedNumber + submissionInterval);

        vm.roll(nextBlockNumber + 1);

        vm.expectEmit(true, true, true, true);
        emit OutputProposed(proposedOutput2, nextOutputIndex, nextBlockNumber, block.timestamp);

        vm.prank(proposer);
        oracle.proposeL2Output(proposedOutput2, nextBlockNumber, 0, 0);
    }

    function setUp() public virtual override {
        super.setUp();
        guardian = makeAddr("guardian");

        // By default the first block has timestamp and number zero, which will cause underflows in the
        // tests, so we'll move forward to these block values.
        initL1Time = startingTimestamp + 1;
        vm.warp(initL1Time);
        vm.roll(startingBlockNumber);
        // Deploy the L2OutputOracle and transfer owernship to the proposer
        oracleImpl = new L2OutputOracle({
            _submissionInterval: submissionInterval,
            _l2BlockTime: l2BlockTime,
            _startingBlockNumber: startingBlockNumber,
            _startingTimestamp: startingTimestamp,
            _proposer: proposer,
            _challenger: owner,
            _finalizationPeriodSeconds: 7 days
        });
        Proxy proxy = new Proxy(multisig);
        vm.prank(multisig);
        proxy.upgradeToAndCall(
            address(oracleImpl),
            abi.encodeCall(L2OutputOracle.initialize, (startingBlockNumber, startingTimestamp))
        );
        oracle = L2OutputOracle(address(proxy));
        vm.label(address(oracle), "L2OutputOracle");

        // Set the L2ToL1MessagePasser at the correct address
        vm.etch(Predeploys.L2_TO_L1_MESSAGE_PASSER, address(new L2ToL1MessagePasser()).code);

        vm.label(Predeploys.L2_TO_L1_MESSAGE_PASSER, "L2ToL1MessagePasser");
    }
}

contract Portal_Initializer is L2OutputOracle_Initializer {
    // Test target
    OptimismPortal internal opImpl;
    OptimismPortal internal op;
    SystemConfig systemConfig;

    event WithdrawalFinalized(bytes32 indexed withdrawalHash, bool success);
    event WithdrawalProven(
        bytes32 indexed withdrawalHash,
        address indexed from,
        address indexed to
    );

    function setUp() public virtual override {
        super.setUp();

        ResourceMetering.ResourceConfig memory config = Constants.DEFAULT_RESOURCE_CONFIG();

        systemConfig = new SystemConfig({
            _owner: address(1),
            _overhead: 0,
            _scalar: 10000,
            _batcherHash: bytes32(0),
            _gasLimit: 30_000_000,
            _unsafeBlockSigner: address(0),
            _config: config
        });

        opImpl = new OptimismPortal({
            _l2Oracle: oracle,
            _guardian: guardian,
            _paused: true,
            _config: systemConfig
        });

        Proxy proxy = new Proxy(multisig);
        vm.prank(multisig);
        proxy.upgradeToAndCall(
            address(opImpl),
            abi.encodeWithSelector(OptimismPortal.initialize.selector, false)
        );
        op = OptimismPortal(payable(address(proxy)));
        vm.label(address(op), "OptimismPortal");
    }
}

contract Messenger_Initializer is Portal_Initializer {
    AddressManager internal addressManager;
    L1CrossDomainMessenger internal L1Messenger;
    L2CrossDomainMessenger internal L2Messenger =
        L2CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER);

    event SentMessage(
        address indexed target,
        address sender,
        bytes message,
        uint256 messageNonce,
        uint256 gasLimit
    );

    event SentMessageExtension1(address indexed sender, uint256 value);

    event MessagePassed(
        uint256 indexed nonce,
        address indexed sender,
        address indexed target,
        uint256 value,
        uint256 gasLimit,
        bytes data,
        bytes32 withdrawalHash
    );

    event RelayedMessage(bytes32 indexed msgHash);
    event FailedRelayedMessage(bytes32 indexed msgHash);

    event TransactionDeposited(
        address indexed from,
        address indexed to,
        uint256 mint,
        uint256 value,
        uint64 gasLimit,
        bool isCreation,
        bytes data
    );

    event WhatHappened(bool success, bytes returndata);

    function setUp() public virtual override {
        super.setUp();

        // Deploy the address manager
        vm.prank(multisig);
        addressManager = new AddressManager();

        // Setup implementation
        L1CrossDomainMessenger L1MessengerImpl = new L1CrossDomainMessenger(op);

        // Setup the address manager and proxy
        vm.prank(multisig);
        addressManager.setAddress("OVM_L1CrossDomainMessenger", address(L1MessengerImpl));
        ResolvedDelegateProxy proxy = new ResolvedDelegateProxy(
            addressManager,
            "OVM_L1CrossDomainMessenger"
        );
        L1Messenger = L1CrossDomainMessenger(address(proxy));
        L1Messenger.initialize();

        vm.etch(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            address(new L2CrossDomainMessenger(address(L1Messenger))).code
        );

        L2Messenger.initialize();

        // Label addresses
        vm.label(address(addressManager), "AddressManager");
        vm.label(address(L1MessengerImpl), "L1CrossDomainMessenger_Impl");
        vm.label(address(L1Messenger), "L1CrossDomainMessenger_Proxy");
        vm.label(Predeploys.LEGACY_ERC20_ETH, "LegacyERC20ETH");
        vm.label(Predeploys.L2_CROSS_DOMAIN_MESSENGER, "L2CrossDomainMessenger");

        vm.label(
            AddressAliasHelper.applyL1ToL2Alias(address(L1Messenger)),
            "L1CrossDomainMessenger_aliased"
        );
    }
}

contract Bridge_Initializer is Messenger_Initializer {
    L1StandardBridge L1Bridge;
    L2StandardBridge L2Bridge;
    OptimismMintableERC20Factory L2TokenFactory;
    OptimismMintableERC20Factory L1TokenFactory;
    ERC20 L1Token;
    ERC20 BadL1Token;
    OptimismMintableERC20 L2Token;
    LegacyMintableERC20 LegacyL2Token;
    ERC20 NativeL2Token;
    ERC20 BadL2Token;
    OptimismMintableERC20 RemoteL1Token;

    event ETHDepositInitiated(address indexed from, address indexed to, uint256 amount, bytes data);

    event ETHWithdrawalFinalized(
        address indexed from,
        address indexed to,
        uint256 amount,
        bytes data
    );

    event ERC20DepositInitiated(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event ERC20WithdrawalFinalized(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event WithdrawalInitiated(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event DepositFinalized(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event DepositFailed(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event ETHBridgeInitiated(address indexed from, address indexed to, uint256 amount, bytes data);

    event ETHBridgeFinalized(address indexed from, address indexed to, uint256 amount, bytes data);

    event ERC20BridgeInitiated(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    event ERC20BridgeFinalized(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes data
    );

    function setUp() public virtual override {
        super.setUp();

        vm.label(Predeploys.L2_STANDARD_BRIDGE, "L2StandardBridge");
        vm.label(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY, "OptimismMintableERC20Factory");

        // Deploy the L1 bridge and initialize it with the address of the
        // L1CrossDomainMessenger
        L1ChugSplashProxy proxy = new L1ChugSplashProxy(multisig);
        vm.mockCall(
            multisig,
            abi.encodeWithSelector(IL1ChugSplashDeployer.isUpgrading.selector),
            abi.encode(true)
        );
        vm.startPrank(multisig);
        proxy.setCode(address(new L1StandardBridge(payable(address(L1Messenger)))).code);
        vm.clearMockedCalls();
        address L1Bridge_Impl = proxy.getImplementation();
        vm.stopPrank();

        L1Bridge = L1StandardBridge(payable(address(proxy)));

        vm.label(address(proxy), "L1StandardBridge_Proxy");
        vm.label(address(L1Bridge_Impl), "L1StandardBridge_Impl");

        // Deploy the L2StandardBridge, move it to the correct predeploy
        // address and then initialize it
        L2StandardBridge l2B = new L2StandardBridge(payable(proxy));
        vm.etch(Predeploys.L2_STANDARD_BRIDGE, address(l2B).code);
        L2Bridge = L2StandardBridge(payable(Predeploys.L2_STANDARD_BRIDGE));

        // Set up the L2 mintable token factory
        OptimismMintableERC20Factory factory = new OptimismMintableERC20Factory(
            Predeploys.L2_STANDARD_BRIDGE
        );
        vm.etch(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY, address(factory).code);
        L2TokenFactory = OptimismMintableERC20Factory(Predeploys.OPTIMISM_MINTABLE_ERC20_FACTORY);

        vm.etch(Predeploys.LEGACY_ERC20_ETH, address(new LegacyERC20ETH()).code);

        L1Token = new ERC20("Native L1 Token", "L1T");

        LegacyL2Token = new LegacyMintableERC20({
            _l2Bridge: address(L2Bridge),
            _l1Token: address(L1Token),
            _name: string.concat("LegacyL2-", L1Token.name()),
            _symbol: string.concat("LegacyL2-", L1Token.symbol())
        });
        vm.label(address(LegacyL2Token), "LegacyMintableERC20");

        // Deploy the L2 ERC20 now
        L2Token = OptimismMintableERC20(
            L2TokenFactory.createStandardL2Token(
                address(L1Token),
                string(abi.encodePacked("L2-", L1Token.name())),
                string(abi.encodePacked("L2-", L1Token.symbol()))
            )
        );

        BadL2Token = OptimismMintableERC20(
            L2TokenFactory.createStandardL2Token(
                address(1),
                string(abi.encodePacked("L2-", L1Token.name())),
                string(abi.encodePacked("L2-", L1Token.symbol()))
            )
        );

        NativeL2Token = new ERC20("Native L2 Token", "L2T");
        L1TokenFactory = new OptimismMintableERC20Factory(address(L1Bridge));

        RemoteL1Token = OptimismMintableERC20(
            L1TokenFactory.createStandardL2Token(
                address(NativeL2Token),
                string(abi.encodePacked("L1-", NativeL2Token.name())),
                string(abi.encodePacked("L1-", NativeL2Token.symbol()))
            )
        );

        BadL1Token = OptimismMintableERC20(
            L1TokenFactory.createStandardL2Token(
                address(1),
                string(abi.encodePacked("L1-", NativeL2Token.name())),
                string(abi.encodePacked("L1-", NativeL2Token.symbol()))
            )
        );
    }
}

contract ERC721Bridge_Initializer is Messenger_Initializer {
    L1ERC721Bridge L1Bridge;
    L2ERC721Bridge L2Bridge;

    function setUp() public virtual override {
        super.setUp();

        // Deploy the L1ERC721Bridge.
        L1Bridge = new L1ERC721Bridge(address(L1Messenger), Predeploys.L2_ERC721_BRIDGE);

        // Deploy the implementation for the L2ERC721Bridge and etch it into the predeploy address.
        vm.etch(
            Predeploys.L2_ERC721_BRIDGE,
            address(new L2ERC721Bridge(Predeploys.L2_CROSS_DOMAIN_MESSENGER, address(L1Bridge)))
                .code
        );

        // Set up a reference to the L2ERC721Bridge.
        L2Bridge = L2ERC721Bridge(Predeploys.L2_ERC721_BRIDGE);

        // Label the L1 and L2 bridges.
        vm.label(address(L1Bridge), "L1ERC721Bridge");
        vm.label(address(L2Bridge), "L2ERC721Bridge");
    }
}

contract FeeVault_Initializer is Bridge_Initializer {
    SequencerFeeVault vault = SequencerFeeVault(payable(Predeploys.SEQUENCER_FEE_WALLET));
    address constant recipient = address(1024);

    event Withdrawal(uint256 value, address to, address from);

    event Withdrawal(
        uint256 value,
        address to,
        address from,
        FeeVault.WithdrawalNetwork withdrawalNetwork
    );
}

contract FFIInterface is Test {
    function getProveWithdrawalTransactionInputs(Types.WithdrawalTransaction memory _tx)
        external
        returns (
            bytes32,
            bytes32,
            bytes32,
            bytes32,
            bytes[] memory
        )
    {
        string[] memory cmds = new string[](8);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "getProveWithdrawalTransactionInputs";
        cmds[2] = vm.toString(_tx.nonce);
        cmds[3] = vm.toString(_tx.sender);
        cmds[4] = vm.toString(_tx.target);
        cmds[5] = vm.toString(_tx.value);
        cmds[6] = vm.toString(_tx.gasLimit);
        cmds[7] = vm.toString(_tx.data);

        bytes memory result = vm.ffi(cmds);
        (
            bytes32 stateRoot,
            bytes32 storageRoot,
            bytes32 outputRoot,
            bytes32 withdrawalHash,
            bytes[] memory withdrawalProof
        ) = abi.decode(result, (bytes32, bytes32, bytes32, bytes32, bytes[]));

        return (stateRoot, storageRoot, outputRoot, withdrawalHash, withdrawalProof);
    }

    function hashCrossDomainMessage(
        uint256 _nonce,
        address _sender,
        address _target,
        uint256 _value,
        uint256 _gasLimit,
        bytes memory _data
    ) external returns (bytes32) {
        string[] memory cmds = new string[](8);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "hashCrossDomainMessage";
        cmds[2] = vm.toString(_nonce);
        cmds[3] = vm.toString(_sender);
        cmds[4] = vm.toString(_target);
        cmds[5] = vm.toString(_value);
        cmds[6] = vm.toString(_gasLimit);
        cmds[7] = vm.toString(_data);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes32));
    }

    function hashWithdrawal(
        uint256 _nonce,
        address _sender,
        address _target,
        uint256 _value,
        uint256 _gasLimit,
        bytes memory _data
    ) external returns (bytes32) {
        string[] memory cmds = new string[](8);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "hashWithdrawal";
        cmds[2] = vm.toString(_nonce);
        cmds[3] = vm.toString(_sender);
        cmds[4] = vm.toString(_target);
        cmds[5] = vm.toString(_value);
        cmds[6] = vm.toString(_gasLimit);
        cmds[7] = vm.toString(_data);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes32));
    }

    function hashOutputRootProof(
        bytes32 _version,
        bytes32 _stateRoot,
        bytes32 _messagePasserStorageRoot,
        bytes32 _latestBlockhash
    ) external returns (bytes32) {
        string[] memory cmds = new string[](6);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "hashOutputRootProof";
        cmds[2] = Strings.toHexString(uint256(_version));
        cmds[3] = Strings.toHexString(uint256(_stateRoot));
        cmds[4] = Strings.toHexString(uint256(_messagePasserStorageRoot));
        cmds[5] = Strings.toHexString(uint256(_latestBlockhash));

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes32));
    }

    function hashDepositTransaction(
        address _from,
        address _to,
        uint256 _mint,
        uint256 _value,
        uint64 _gas,
        bytes memory _data,
        uint64 _logIndex
    ) external returns (bytes32) {
        string[] memory cmds = new string[](10);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "hashDepositTransaction";
        cmds[2] = "0x0000000000000000000000000000000000000000000000000000000000000000";
        cmds[3] = vm.toString(_logIndex);
        cmds[4] = vm.toString(_from);
        cmds[5] = vm.toString(_to);
        cmds[6] = vm.toString(_mint);
        cmds[7] = vm.toString(_value);
        cmds[8] = vm.toString(_gas);
        cmds[9] = vm.toString(_data);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes32));
    }

    function encodeDepositTransaction(Types.UserDepositTransaction calldata txn)
        external
        returns (bytes memory)
    {
        string[] memory cmds = new string[](11);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "encodeDepositTransaction";
        cmds[2] = vm.toString(txn.from);
        cmds[3] = vm.toString(txn.to);
        cmds[4] = vm.toString(txn.value);
        cmds[5] = vm.toString(txn.mint);
        cmds[6] = vm.toString(txn.gasLimit);
        cmds[7] = vm.toString(txn.isCreation);
        cmds[8] = vm.toString(txn.data);
        cmds[9] = vm.toString(txn.l1BlockHash);
        cmds[10] = vm.toString(txn.logIndex);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes));
    }

    function encodeCrossDomainMessage(
        uint256 _nonce,
        address _sender,
        address _target,
        uint256 _value,
        uint256 _gasLimit,
        bytes memory _data
    ) external returns (bytes memory) {
        string[] memory cmds = new string[](8);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "encodeCrossDomainMessage";
        cmds[2] = vm.toString(_nonce);
        cmds[3] = vm.toString(_sender);
        cmds[4] = vm.toString(_target);
        cmds[5] = vm.toString(_value);
        cmds[6] = vm.toString(_gasLimit);
        cmds[7] = vm.toString(_data);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (bytes));
    }

    function decodeVersionedNonce(uint256 nonce) external returns (uint256, uint256) {
        string[] memory cmds = new string[](3);
        cmds[0] = "scripts/differential-testing/differential-testing";
        cmds[1] = "decodeVersionedNonce";
        cmds[2] = vm.toString(nonce);

        bytes memory result = vm.ffi(cmds);
        return abi.decode(result, (uint256, uint256));
    }

    function getMerkleTrieFuzzCase(string memory variant)
        external
        returns (
            bytes32,
            bytes memory,
            bytes memory,
            bytes[] memory
        )
    {
        string[] memory cmds = new string[](5);
        cmds[0] = "./test-case-generator/fuzz";
        cmds[1] = "-m";
        cmds[2] = "trie";
        cmds[3] = "-v";
        cmds[4] = variant;

        return abi.decode(vm.ffi(cmds), (bytes32, bytes, bytes, bytes[]));
    }
}

// Used for testing a future upgrade beyond the current implementations.
// We include some variables so that we can sanity check accessing storage values after an upgrade.
contract NextImpl is Initializable {
    // Initializable occupies the zero-th slot.
    bytes32 slot1;
    bytes32[19] __gap;
    bytes32 slot21;
    bytes32 public constant slot21Init = bytes32(hex"1337");

    function initialize() public reinitializer(2) {
        // Slot21 is unused by an of our upgradeable contracts.
        // This is used to verify that we can access this value after an upgrade.
        slot21 = slot21Init;
    }
}

contract Reverter {
    fallback() external {
        revert();
    }
}

// Useful for testing reentrancy guards
contract CallerCaller {
    event WhatHappened(bool success, bytes returndata);

    fallback() external {
        (bool success, bytes memory returndata) = msg.sender.call(msg.data);
        emit WhatHappened(success, returndata);
        assembly {
            switch success
            case 0 {
                revert(add(returndata, 0x20), mload(returndata))
            }
            default {
                return(add(returndata, 0x20), mload(returndata))
            }
        }
    }
}

// Used for testing the `CrossDomainMessenger`'s per-message reentrancy guard.
contract ConfigurableCaller {
    bool doRevert = true;
    address target;
    bytes payload;

    event WhatHappened(bool success, bytes returndata);

    /**
     * @notice Call the configured target with the configured payload OR revert.
     */
    function call() external {
        if (doRevert) {
            revert("ConfigurableCaller: revert");
        } else {
            (bool success, bytes memory returndata) = address(target).call(payload);
            emit WhatHappened(success, returndata);
            assembly {
                switch success
                case 0 {
                    revert(add(returndata, 0x20), mload(returndata))
                }
                default {
                    return(add(returndata, 0x20), mload(returndata))
                }
            }
        }
    }

    /**
     * @notice Set whether or not to have `call` revert.
     */
    function setDoRevert(bool _doRevert) external {
        doRevert = _doRevert;
    }

    /**
     * @notice Set the target for the call made in `call`.
     */
    function setTarget(address _target) external {
        target = _target;
    }

    /**
     * @notice Set the payload for the call made in `call`.
     */
    function setPayload(bytes calldata _payload) external {
        payload = _payload;
    }

    /**
     * @notice Fallback function that reverts if `doRevert` is true.
     *         Otherwise, it does nothing.
     */
    fallback() external {
        if (doRevert) {
            revert("ConfigurableCaller: revert");
        }
    }
}
