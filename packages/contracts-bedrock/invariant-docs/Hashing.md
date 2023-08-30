# `Hashing` Invariants

## `hashCrossDomainMessage` reverts if `version` is > `1`.
**Test:** [`Hashing.t.sol#L141`](../contracts/test/invariants/Hashing.t.sol#L141)

The `hashCrossDomainMessage` function should always revert if the `version` passed is > `1`. 


## `version` = `0`: `hashCrossDomainMessage` and `hashCrossDomainMessageV0` are equivalent.
**Test:** [`Hashing.t.sol#L153`](../contracts/test/invariants/Hashing.t.sol#L153)

If the version passed is 0, `hashCrossDomainMessage` and `hashCrossDomainMessageV0` should be equivalent. 


## `version` = `1`: `hashCrossDomainMessage` and `hashCrossDomainMessageV1` are equivalent.
**Test:** [`Hashing.t.sol#L166`](../contracts/test/invariants/Hashing.t.sol#L166)

If the version passed is 1, `hashCrossDomainMessage` and `hashCrossDomainMessageV1` should be equivalent. 
