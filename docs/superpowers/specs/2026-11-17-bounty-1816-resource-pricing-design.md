# Bounty 1816: charge resource gas at the transaction price

The user approved actual-price conversion with an enforced minimum transaction
gas price of **0.05 gwei**. Work uses `fix-1816-bugbounty`, based on
`develop-bug-bounty` at `ea7b00598db0bd6bf5ba4c8ddd1d1252080c5fb4`.

## Behavior and scope

Replace ResourceMetering's fixed 3 gwei divisor with
`Math.ceilDiv(resourceCost, tx.gasprice)`. Before updating metering state,
reject transaction prices below 0.05 gwei with `GasPriceTooLow()`. This includes
zero-price transactions even when the requested resource amount is zero.
Define the minimum as an internal constant so storage layout is preserved.

Preserve the existing execution-gas credit, resource cap, EIP-1559-style fee
curve, deposit entry points, deposit values and events. The shared base class
applies the correction to both portal implementations. Normal low-price
deposits require more execution gas; at a resource price of 1 gwei and a
transaction price of 0.05 gwei, 100,000 resource units cost about 2,000,000
execution gas. An 80,000,000-unit purchase cannot fit in 55,000,000 execution
gas at that price and must roll back all resource accounting.

Keeping the historical divisor and requiring 3 gwei transactions would reject
ordinary lower-price deposits. Charging a separate explicit payment would
require a broader change to the portal's native and ERC20 deposit interfaces.
Neither alternative is part of the approved fix.

## Validation

First run the economic regression against the unchanged production contract:
a 0.05 gwei deposit must pay at least its nominal resource cost, and a full-cap
purchase with a 55,000,000 execution-gas budget must fail without consuming
capacity. Both assertions must fail on the original implementation.

Cover the 0.05 gwei boundary, zero and sub-floor prices, other transaction
prices, non-divisible resource costs, negligible requested amounts, same-block
capacity, next-block recovery, alternating blocks, nonempty intervening blocks
and adjacent-block fee escalation. Check portal deposits preserve native value
and event data and revert atomically for invalid prices or insufficient gas.

Run the existing resource-metering tests and relevant portal coverage. Set the
test fixture default transaction gas price to the historical 3 gwei assumption;
regressions explicitly override it. If full-project compilation remains
unavailable, use an isolated Foundry project referencing the actual repository
source and test files and report exactly which coverage ran.

## Constraints

- Enforce `tx.gasprice >= 0.05 gwei` in the contract.
- Use actual transaction gas price and round gas requirements upward.
- Preserve storage layout, deposit interfaces and deposit accounting.
- Preserve the resource cap and fee adjustment algorithm.
- Do not change dependencies or unrelated code.
- Preserve all pre-existing untracked documents.

## Approval

The user selected “按实际价格换算，强制最低 0.05 gwei（推荐）” after reviewing
the proposed calculation, floor and compatibility impact in this chat.
