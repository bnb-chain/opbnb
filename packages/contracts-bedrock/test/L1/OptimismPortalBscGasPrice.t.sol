// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { MockERC20 } from "solmate/test/utils/mocks/MockERC20.sol";
import { OptimismPortal } from "src/L1/OptimismPortal.sol";
import { OptimismPortal2 } from "src/L1/OptimismPortal2.sol";
import { L2OutputOracle } from "src/L1/L2OutputOracle.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { Proxy } from "src/universal/Proxy.sol";
import { DisputeGameFactory } from "src/dispute/DisputeGameFactory.sol";
import { GameType } from "src/dispute/lib/Types.sol";

/// @notice Exercises actual portal and SystemConfig proxies without the full L2 genesis fixture.
contract OptimismPortalBscGasPrice_Test is Test {
    event TransactionDeposited(address indexed from, address indexed to, uint256 indexed version, bytes opaqueData);

    uint256 internal constant REFERENCE_GAS_PRICE = 0.05 gwei;
    OptimismPortal internal portal;
    address internal depositor;
    address internal recipient;

    function setUp() public {
        vm.fee(0);
        vm.txGasPrice(REFERENCE_GAS_PRICE);
        depositor = makeAddr("depositor");
        recipient = makeAddr("recipient");
        vm.deal(depositor, 1 ether);
        portal = _deployPortal(address(0));
    }

    function test_depositTransaction_preservesValueAndEvent_succeeds() external {
        vm.expectEmit(true, true, true, true, address(portal));
        emit TransactionDeposited(
            depositor,
            recipient,
            0,
            abi.encodePacked(uint256(1 ether), uint256(0.3 ether), uint64(100_000), false, hex"1234")
        );
        vm.prank(depositor, depositor);
        uint256 gasBefore = gasleft();
        portal.depositTransaction{ value: 1 ether }(recipient, 0.3 ether, 100_000, false, hex"1234");
        uint256 gasUsed = gasBefore - gasleft();

        assertGe(gasUsed * tx.gasprice, 100_000 * 1 gwei);
        assertEq(address(portal).balance, 1 ether);
        assertEq(portal.balance(), 1 ether);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 100_000);
    }

    function test_receive_preservesValueAndEvent_succeeds() external {
        vm.txGasPrice(0.01 gwei);
        vm.expectEmit(true, true, true, true, address(portal));
        emit TransactionDeposited(
            depositor,
            depositor,
            0,
            abi.encodePacked(uint256(1 ether), uint256(1 ether), uint64(100_000), false, bytes(""))
        );
        vm.prank(depositor, depositor);
        (bool success,) = address(portal).call{ value: 1 ether }("");
        assertTrue(success);
        assertEq(portal.balance(), 1 ether);
    }

    function test_depositTransaction_belowHistoricalGasPrice_succeeds() external {
        vm.txGasPrice(0.01 gwei);
        vm.prank(depositor, depositor);
        uint256 gasBefore = gasleft();
        portal.depositTransaction{ value: 1 ether }(recipient, 0, 100_000, false, "");
        uint256 gasUsed = gasBefore - gasleft();

        assertGe(gasUsed * tx.gasprice, 100_000 * 1 gwei);
        assertEq(address(portal).balance, 1 ether);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 100_000);
    }

    function test_depositTransaction_zeroGasPrice_reverts() external {
        vm.txGasPrice(0);
        vm.expectRevert(ResourceMetering.GasPriceTooLow.selector);
        vm.prank(depositor, depositor);
        portal.depositTransaction{ value: 1 ether }(recipient, 0, 100_000, false, "");

        assertEq(address(portal).balance, 0);
        assertEq(depositor.balance, 1 ether);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 0);
    }

    function test_depositTransaction_tinyGasPriceExceedsBudget_reverts() external {
        vm.txGasPrice(1);
        vm.prank(depositor, depositor);
        (bool success, bytes memory reason) = address(portal).call{ gas: 1_000_000, value: 1 ether }(
            abi.encodeCall(portal.depositTransaction, (recipient, 0, 100_000, false, bytes("")))
        );
        assertFalse(success);
        assertEq(reason.length, 0, "positive prices must fail from gas exhaustion, not a price check");
        assertEq(address(portal).balance, 0);
        assertEq(depositor.balance, 1 ether);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 0);
    }

    function test_depositTransaction_zeroValue_succeeds() external {
        vm.prank(depositor, depositor);
        portal.depositTransaction(recipient, 0, 100_000, false, "");
        assertEq(portal.balance(), 0);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 100_000);
    }

    function test_depositTransaction_fullCapExceedsBlockGas_reverts() external {
        vm.prank(depositor, depositor);
        (bool success,) = address(portal).call{ gas: 55_000_000, value: 1 ether }(
            abi.encodeCall(portal.depositTransaction, (recipient, 0, 80_000_000, false, bytes("")))
        );
        assertFalse(success);
        assertEq(address(portal).balance, 0);
        assertEq(depositor.balance, 1 ether);
        (, uint64 boughtGas,) = portal.params();
        assertEq(boughtGas, 0);

        vm.prank(depositor, depositor);
        portal.depositTransaction(recipient, 0, 100_000, false, "");
    }

    function test_depositERC20Transaction_preservesTokensAndEvent_succeeds() external {
        vm.txGasPrice(0.01 gwei);
        MockERC20 token = new MockERC20("Test", "TST", 18);
        OptimismPortal customPortal = _deployPortal(address(token));
        token.mint(depositor, 1 ether);
        vm.prank(depositor);
        token.approve(address(customPortal), 1 ether);

        vm.expectEmit(true, true, true, true, address(customPortal));
        emit TransactionDeposited(
            depositor,
            recipient,
            0,
            abi.encodePacked(uint256(1 ether), uint256(0.3 ether), uint64(100_000), false, bytes(""))
        );
        vm.prank(depositor, depositor);
        uint256 gasBefore = gasleft();
        customPortal.depositERC20Transaction(recipient, 1 ether, 0.3 ether, 100_000, false, "");
        uint256 gasUsed = gasBefore - gasleft();

        assertGe(gasUsed * tx.gasprice, 100_000 * 1 gwei);
        assertEq(token.balanceOf(address(customPortal)), 1 ether);
        assertEq(token.balanceOf(depositor), 0);
        assertEq(customPortal.balance(), 1 ether);
    }

    function test_depositERC20Transaction_zeroGasPrice_reverts() external {
        MockERC20 token = new MockERC20("Test", "TST", 18);
        OptimismPortal customPortal = _deployPortal(address(token));
        token.mint(depositor, 1 ether);
        vm.prank(depositor);
        token.approve(address(customPortal), 1 ether);
        (, uint64 boughtGasBefore,) = customPortal.params();

        vm.txGasPrice(0);
        vm.expectRevert(ResourceMetering.GasPriceTooLow.selector);
        vm.prank(depositor, depositor);
        customPortal.depositERC20Transaction(recipient, 1 ether, 0, 100_000, false, "");

        assertEq(token.balanceOf(depositor), 1 ether);
        assertEq(token.balanceOf(address(customPortal)), 0);
        assertEq(customPortal.balance(), 0);
        (, uint64 boughtGasAfter,) = customPortal.params();
        assertEq(boughtGasAfter, boughtGasBefore);
    }

    function test_portal2_preservesValueAndEvent_succeeds() external {
        vm.txGasPrice(0.01 gwei);
        OptimismPortal2 portal2 = _deployPortal2();
        vm.expectEmit(true, true, true, true, address(portal2));
        emit TransactionDeposited(
            depositor,
            recipient,
            0,
            abi.encodePacked(uint256(1 ether), uint256(0.3 ether), uint64(100_000), false, bytes(""))
        );
        vm.prank(depositor, depositor);
        uint256 gasBefore = gasleft();
        portal2.depositTransaction{ value: 1 ether }(recipient, 0.3 ether, 100_000, false, "");
        uint256 gasUsed = gasBefore - gasleft();
        assertGe(gasUsed * tx.gasprice, 100_000 * 1 gwei);
        assertEq(address(portal2).balance, 1 ether);
        (, uint64 boughtGas,) = portal2.params();
        assertEq(boughtGas, 100_000);
    }

    function test_portal2_zeroGasPrice_reverts() external {
        OptimismPortal2 portal2 = _deployPortal2();
        vm.txGasPrice(0);
        vm.expectRevert(ResourceMetering.GasPriceTooLow.selector);
        vm.prank(depositor, depositor);
        portal2.depositTransaction{ value: 1 ether }(recipient, 0, 100_000, false, "");
        assertEq(address(portal2).balance, 0);
        assertEq(depositor.balance, 1 ether);
        (, uint64 boughtGas,) = portal2.params();
        assertEq(boughtGas, 0);
    }

    function test_portal2_fullCapExceedsBlockGas_reverts() external {
        OptimismPortal2 portal2 = _deployPortal2();
        vm.prank(depositor, depositor);
        (bool success,) = address(portal2).call{ gas: 55_000_000, value: 1 ether }(
            abi.encodeCall(portal2.depositTransaction, (recipient, 0, 80_000_000, false, bytes("")))
        );
        assertFalse(success);
        assertEq(address(portal2).balance, 0);
        assertEq(depositor.balance, 1 ether);
        (, uint64 boughtGas,) = portal2.params();
        assertEq(boughtGas, 0);
    }

    function _deployPortal2() internal returns (OptimismPortal2 portal_) {
        Proxy proxy = new Proxy(address(this));
        proxy.upgradeToAndCall(
            address(new OptimismPortal2(0, 0)),
            abi.encodeCall(
                OptimismPortal2.initialize,
                (DisputeGameFactory(address(0)), portal.systemConfig(), SuperchainConfig(address(0)), GameType.wrap(0))
            )
        );
        portal_ = OptimismPortal2(payable(address(proxy)));
    }

    function _deployPortal(address _gasToken) internal returns (OptimismPortal portal_) {
        Proxy configProxy = new Proxy(address(this));
        Proxy portalProxy = new Proxy(address(this));
        SystemConfig config = SystemConfig(address(configProxy));
        portal_ = OptimismPortal(payable(address(portalProxy)));
        portalProxy.upgradeToAndCall(
            address(new OptimismPortal()),
            abi.encodeCall(
                OptimismPortal.initialize, (L2OutputOracle(address(0)), config, SuperchainConfig(address(0)))
            )
        );

        SystemConfig.Addresses memory addresses;
        addresses.optimismPortal = address(portal_);
        addresses.gasPayingToken = _gasToken;
        configProxy.upgradeToAndCall(
            address(new SystemConfig()),
            abi.encodeCall(
                SystemConfig.initialize,
                (
                    address(this),
                    0,
                    0,
                    bytes32(0),
                    100_000_000,
                    address(0),
                    ResourceMetering.ResourceConfig({
                        maxResourceLimit: 80_000_000,
                        elasticityMultiplier: 2,
                        baseFeeMaxChangeDenominator: 8,
                        minimumBaseFee: 1 gwei,
                        systemTxMaxGas: 1_000_000,
                        maximumBaseFee: type(uint128).max
                    }),
                    address(0),
                    addresses
                )
            )
        );
    }
}
