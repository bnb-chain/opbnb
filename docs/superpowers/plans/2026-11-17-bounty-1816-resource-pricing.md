# Bounty 1816 Resource Pricing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [x]`) syntax for tracking.

**Goal:** Remove the fixed-divisor resource discount and enforce the approved 0.05 gwei minimum gas price.

**Architecture:** Retain ResourceMetering's gas-burn mechanism and execution-gas credit. Convert the nominal resource cost with OpenZeppelin's existing `Math.ceilDiv` and actual `tx.gasprice`; reject prices below the minimum before metering state updates. Use the real production base class and portals in tests.

**Tech Stack:** Solidity 0.8.15, Foundry, forge-std, existing OpenZeppelin Math library.

## Global Constraints

- Enforce `tx.gasprice >= 0.05 gwei` in the contract.
- Use actual transaction gas price and round gas requirements upward.
- Preserve storage layout, deposit interfaces and deposit accounting.
- Preserve the resource cap and fee adjustment algorithm.
- Do not change dependencies or unrelated code.
- Preserve all pre-existing untracked documents.

---

### Task 1: Correct resource pricing and verify deposit compatibility

**Files:**
- Modify: `packages/contracts-bedrock/src/L1/ResourceMetering.sol`
- Modify: `packages/contracts-bedrock/test/setup/CommonTest.sol`
- Modify: `packages/contracts-bedrock/test/L1/ResourceMetering.t.sol`
- Create: `packages/contracts-bedrock/test/L1/ResourceMeteringBscGasPrice.t.sol`
- Create: `packages/contracts-bedrock/test/L1/OptimismPortalBscGasPrice.t.sol`
- Update: `packages/contracts-bedrock/snapshots/abi/OptimismPortal.json`, `OptimismPortal2.json`, `OptimismPortalInterop.json` (new error only)

**Interfaces:**
- Consumes: `metered(uint64)`, `_metered(uint64,uint256)`, `params()`, and the unchanged portal deposit interfaces.
- Produces: `GasPriceTooLow()` and corrected resource pricing; no new callable functions or storage fields.

- [x] Add economic regressions using the existing `MeterUser` test harness with the report's 80,000,000-unit resource cap and elasticity 2. Measure `gasBefore - gasleft()` around the actual call and assert `gasUsed * tx.gasprice >= amount * resourceBaseFee`. For the attack, call the meter with 55,000,000 gas, assert failure and unchanged bought gas, then verify a normal deposit succeeds in the same block.
- [x] Run `forge test --offline --match-contract ResourceMeteringBscGasPrice_Test -vv`. Observe nominal-price and bounded-full-cap assertions fail on the original contract. If whole-project solc processes are killed, use a temporary Foundry root with empty source/script directories, symlink the repository tests, and map `src/` and existing libraries back to the repository. Use solc 0.8.15 and London for these contracts.
- [x] Add the minimal production correction:

```solidity
error GasPriceTooLow();
uint256 internal constant MINIMUM_GAS_PRICE = 0.05 gwei;
// At the start of _metered:
if (tx.gasprice < MINIMUM_GAS_PRICE) revert GasPriceTooLow();
// In place of the fixed-divisor conversion:
uint256 gasCost = Math.ceilDiv(resourceCost, tx.gasprice);
```

- [x] Add `vm.txGasPrice(3 gwei)` to CommonTest and both existing ResourceMetering test setups to make the historical price assumption explicit. Override gas price in every new pricing regression. Keep global script gas-price defaults unchanged.
- [x] Cover price floor and zero amount, 0.05/0.1/1/3/30 gwei prices, non-divisible costs, same-block caps, next-block recovery, alternating blocks and adjacent-block fee escalation. Verify native and custom-token deposits, receive(), event data and rollback through production portal code.
- [x] Run the new regressions, existing ResourceMetering tests, and relevant portal tests. Check `forge fmt --check` for changed Solidity files and `git diff --check`; compare compiled storage layouts against baseline. Review the final diff for changes to the fee curve, interfaces or unrelated documents.
- [x] Commit only this fix's source, tests, ABI snapshots and its own design/plan documents with `fix(contracts): price deposit resources at actual transaction gas price`.

## Execution

The user authorized implementation by selecting the recommended design. This
single task has tightly coupled red/green steps and will be executed inline;
no further approval is needed for implementation or verification.

## Validation results

- On the unchanged implementation, 14 new regression cases failed for the
  expected economic or missing-price-check assertions; 8 controls passed.
- After the fix, all 25 new cases and all 12 original ResourceMetering cases
  passed. The original CSV artifact test also passed with a 100,000,000,000
  gas test budget (it consumes approximately 4.81 billion gas).
- Both actual portal implementations were deployed behind production Proxy
  contracts with a real SystemConfig proxy. Native deposits, ERC20 deposits,
  receive(), event data, balances and rollback were verified.
- ResourceMetering and OptimismPortal compiled storage layouts match baseline.
  All three portal ABI snapshots differ only by GasPriceTooLow().
- forge fmt --check and git diff --check passed.

The installed ARM64 solc 0.8.15 was repeatedly terminated while compiling
forge-std/Test.sol. Validation used the official x86_64 solc 0.8.15 release
(e14f2714), Foundry 1.3.6, and a temporary London-targeted project with the repository's optimizer settings and bytecode hash disabled, referencing
actual repository sources and tests. No repository dependencies or global
compiler installation were changed. The full contracts suite was not run.

Commands in that temporary root:

```sh
forge test --offline --use ./solc-macos-official --no-match-contract ArtifactResourceMetering_Test -vv
forge test --offline --use ./solc-macos-official --match-contract ArtifactResourceMetering_Test --gas-limit 100000000000 -vv
```
