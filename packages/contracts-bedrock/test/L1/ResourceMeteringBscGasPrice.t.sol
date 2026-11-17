// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { MeterUser } from "test/L1/ResourceMetering.t.sol";

/// @notice Regressions for SRC-1816 using the reported BSC resource configuration.
contract ResourceMeteringBscGasPrice_Test is Test {
    uint256 internal constant REFERENCE_GAS_PRICE = 0.05 gwei;
    uint32 internal constant MAX_RESOURCE_LIMIT = 80_000_000;
    MeterUser internal meter;

    function setUp() public {
        vm.fee(0);
        vm.txGasPrice(3 gwei);
        meter = new MeterUser();
        ResourceMetering.ResourceConfig memory config = meter.resourceConfig();
        config.maxResourceLimit = MAX_RESOURCE_LIMIT;
        config.elasticityMultiplier = 2;
        meter.setParams(config);
    }

    function test_meter_lowGasPricePaysNominalCost_succeeds() external {
        _assertCharged(100_000, REFERENCE_GAS_PRICE);
    }

    function test_meter_transactionPricesPayNominalCost_succeeds() external {
        _assertCharged(100_000, 0.01 gwei);
        _assertCharged(1_000_000, 0.1 gwei);
        _assertCharged(1_000_000, 1 gwei);
        _assertCharged(1_000_000, 3 gwei);
        _assertCharged(1_000_000, 30 gwei);
    }

    function test_meter_fractionalGasCost_succeeds() external {
        _assertCharged(100_001, 3 gwei);
    }

    function test_meter_executionGasCoversCost_succeeds() external {
        _assertCharged(1, 30 gwei);
    }

    function test_meter_zeroGasPrice_reverts() external {
        vm.txGasPrice(0);
        vm.expectRevert(ResourceMetering.GasPriceTooLow.selector);
        meter.use(100_000);
    }

    function test_meter_belowHistoricalGasPrice_succeeds() external {
        _assertCharged(100_000, 0.01 gwei);
        _assertCharged(100_000, REFERENCE_GAS_PRICE - 1);
    }

    function test_meter_tinyGasPriceExceedsBudget_reverts() external {
        meter.use(100_000);
        (uint128 baseFeeBefore, uint64 boughtGasBefore, uint64 blockNumBefore) = meter.params();
        vm.roll(block.number + 1);
        vm.txGasPrice(1);

        (bool success, bytes memory reason) =
            address(meter).call{ gas: 1_000_000 }(abi.encodeCall(meter.use, (uint64(100_000))));
        assertFalse(success);
        assertEq(reason.length, 0, "positive prices must fail from gas exhaustion, not a price check");
        (uint128 baseFeeAfter, uint64 boughtGasAfter, uint64 blockNumAfter) = meter.params();
        assertEq(baseFeeAfter, baseFeeBefore);
        assertEq(boughtGasAfter, boughtGasBefore);
        assertEq(blockNumAfter, blockNumBefore);
    }

    function test_meter_oneWeiGasPriceZeroAmount_succeeds() external {
        vm.txGasPrice(1);
        meter.use(0);
        (, uint64 boughtGas,) = meter.params();
        assertEq(boughtGas, 0);
    }

    function test_meter_zeroAmountWithZeroGasPrice_reverts() external {
        vm.txGasPrice(0);
        vm.expectRevert(ResourceMetering.GasPriceTooLow.selector);
        meter.use(0);
    }

    function test_meter_zeroAmountAtLowGasPrice_succeeds() external {
        vm.txGasPrice(REFERENCE_GAS_PRICE);
        meter.use(0);
        (, uint64 boughtGas,) = meter.params();
        assertEq(boughtGas, 0);
    }

    function test_meter_fullCapAtLowPriceExceedsBlockGas_reverts() external {
        vm.txGasPrice(REFERENCE_GAS_PRICE);
        (bool success,) = address(meter).call{ gas: 55_000_000 }(abi.encodeCall(meter.use, (MAX_RESOURCE_LIMIT)));
        assertFalse(success, "full-cap purchase must exceed the execution gas budget");

        (uint128 baseFee, uint64 boughtGas, uint64 blockNum) = meter.params();
        assertEq(baseFee, 1 gwei);
        assertEq(boughtGas, 0, "failed purchase must not consume deposit capacity");
        assertEq(blockNum, block.number);
        _assertCharged(100_000, REFERENCE_GAS_PRICE);
    }

    function test_meter_sameBlockResourceCap_reverts() external {
        _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
        vm.expectRevert(ResourceMetering.OutOfGas.selector);
        meter.use(1);
        (, uint64 boughtGas,) = meter.params();
        assertEq(boughtGas, MAX_RESOURCE_LIMIT);
    }

    function test_meter_nextBlockRestoresCapacity_succeeds() external {
        _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
        vm.roll(block.number + 1);
        _assertCharged(100_000, REFERENCE_GAS_PRICE);
        (, uint64 boughtGas,) = meter.params();
        assertEq(boughtGas, 100_000);
    }

    function test_meter_alternatingBlocksPayNominalCost_succeeds() external {
        for (uint256 i; i < 3; i++) {
            _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
            (uint128 baseFee, uint64 boughtGas,) = meter.params();
            assertEq(baseFee, 1 gwei);
            assertEq(boughtGas, MAX_RESOURCE_LIMIT);
            vm.roll(block.number + 2);
        }
    }

    function test_meter_adjacentBlocksRaiseBaseFee_succeeds() external {
        _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
        vm.roll(block.number + 1);
        _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
        (uint128 baseFee,,) = meter.params();
        assertEq(baseFee, 1.125 gwei);
    }

    function test_meter_nonemptyInterveningBlockRetainsFeeIncrease_succeeds() external {
        _assertCharged(MAX_RESOURCE_LIMIT, 3 gwei);
        vm.roll(block.number + 1);
        _assertCharged(5_000_000, 3 gwei);
        vm.roll(block.number + 1);
        _assertCharged(100_000, REFERENCE_GAS_PRICE);
        (uint128 baseFee,,) = meter.params();
        assertEq(baseFee, 1_001_953_125);
    }

    function testFuzz_meter_actualPricePaysNominalCost_succeeds(uint64 _amount, uint64 _gasPrice) external {
        _amount = uint64(bound(_amount, 21_000, 200_000));
        // Keep successful fuzz purchases within the test's execution gas budget.
        _gasPrice = uint64(bound(_gasPrice, 0.01 gwei, 30 gwei));
        _assertCharged(_amount, _gasPrice);
    }

    function _assertCharged(uint64 _amount, uint256 _gasPrice) internal {
        vm.txGasPrice(_gasPrice);
        uint256 gasBefore = gasleft();
        meter.use(_amount);
        uint256 gasUsed = gasBefore - gasleft();

        (uint128 baseFee,,) = meter.params();
        uint256 nominalCost = uint256(_amount) * uint256(baseFee);
        assertGe(gasUsed * _gasPrice, nominalCost, "metered execution gas must cover the resource cost");
        // Allow for call overhead and the execution-gas credit when no burn is necessary.
        assertLe(gasUsed, (nominalCost + _gasPrice - 1) / _gasPrice + 30_000, "excessive resource gas burn");
    }
}
