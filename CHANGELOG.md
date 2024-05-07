# Changelog

## v0.3.3

This is a minor release and upgrading is optional.

### User Facing Changes

To simplify the process of starting op-node for users, default configurations for opBNB mainnet and testnet have been added. Users can now select the network configuration by setting `--network=opBNBMainnet` or `--network=opBNBTestnet`. Check details in PR #179.

### Partial Changelog

- #179: feature(op-node): simplify op-node start

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.3.3

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.3.2...v0.3.3

## v0.3.2

This is a hardfork release for the opBNB Mainnet called Snow. It will be activated at April 15, 2024, at 6 AM UTC.
All op-node of mainnet nodes must upgrade to this release before the hardfork.

**If you are on v0.3.1, you need to upgrade to v0.3.2 by April 9, 2024, at 6 AM UTC.**

### User Facing Changes

- The L1 gas price of all L2 transactions will be optimized after the snow hardfork. The price will be calculated based on the median of the last 21 blocks' gas prices on BSC. The L1 gas price for the opBNB Mainnet is expected to be decreased to 1 Gwei after the hardfork. And it will adjust automatically if the gas price on BSC changes.

### Partial Changelog

- #169: feat: optimize l1 gas price calculation after snow hardfork

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.3.2

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.3.0...v0.3.2

## v0.3.1 [Deprecated]

This is a hardfork release for the opBNB Mainnet called Snow. It will be activated at April 9, 2024, at 6 AM UTC.
All mainnet nodes must upgrade to this release before the hardfork.

**Please notice that this version is deprecated, and you need to upgrade to v0.3.2 by April 9, 2024, at 6 AM UTC if you are on this version.**

### User Facing Changes

- The L1 gas price of all L2 transactions will be optimized after the snow hardfork. The price will be calculated based on the median of the last 21 blocks' gas prices on BSC. The L1 gas price for the opBNB Mainnet is expected to be decreased to 1 Gwei after the hardfork. And it will adjust automatically if the gas price on BSC changes.

### Partial Changelog

- #169: feat: optimize l1 gas price calculation after snow hardfork

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.3.1
- ghcr.io/bnb-chain/op-batcher:v0.3.1
- ghcr.io/bnb-chain/op-proposer:v0.3.1

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.3.0...v0.3.1

## v0.3.0

This is a recommended release for op-node. This release brings in upstream updates, see https://github.com/bnb-chain/opbnb/pull/121 for the contents. This is also a ready release for the next planed fork, which will bring in canyon fork from upstream as well.

### User Facing Changes

- The default address for the metrics, rpc, and pprof servers will be changing from 0.0.0.0 to 127.0.0.1, and pprof server will change to use 6060 port by default(you may encounter panic due to port conflict, if so you could change it to another value via `--pprof.port`).
- op-node enable p2p score feature, and with `--p2p.scoring=light`, `--p2p.ban.peers=true` by default(to keep previous config, you could set `--p2p.scoring=none`, `--p2p.ban.peers=false`).

### Partial Changelog

- #115: fix(op-node): pre-fetching handle L1 reOrg
- #117: fix(op-node): pre-fetching handle L1 reOrg round 2
- #118: fix(op-node/op-batcher/op-proposer): the fallback client should always try recover
- #121: Merge upstream v1.3.0
- #127: fix(op-node): fix basefee when start new chain with fermat

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.3.0
- ghcr.io/bnb-chain/op-batcher:v0.3.0
- ghcr.io/bnb-chain/op-proposer:v0.3.0

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.2.4...v0.3.0

## v0.2.4

This is a minor release and upgrading is optional.

### User Facing Changes

- The op-node is pre-configured with bootnodes for opBNB mainnet and testnet. By default, the op-node will use these pre-configured bootnodes for both networks. If you prefer to use your own bootnodes, you can configure the `p2p.bootnodes` parameter in the command line flag.(#89)

### Partial Changelog

- #87: optimize(op-node): make block produce stable when L1 latency unstable
- #89: feat(op-node): add opBNB bootnodes
- #94: fix(op-node/op-batcher): fallbackClient should ignore ethereum.NotFound error
- #100: feature(op-node): pre-fetch receipts concurrently
- #101: optimize(op-node): continue optimizing sequencer step schedule
- #104: feat(op-node): pre-fetch receipts concurrently round 2
- #106: optimize: extended expire time for sequencer block broadcasting
- #108: optimize(op-node): increase catching up speed when sequencer lagging
- #109: feat(op-batcher/op-proposer): add InstrumentedClient
- #111: fix(op-node): remove 3s stepCtx for sequencer

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.2.4
- ghcr.io/bnb-chain/op-batcher:v0.2.4
- ghcr.io/bnb-chain/op-proposer:v0.2.4

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.2.3...v0.2.4

## v0.2.3

This is a minor release and upgrading is optional.

### User Facing Changes

NA

### Partial Changelog

- #88: fix: change fallback client threshold to 10 from 20
- #90: fix(op-node): not do sequence action instantly

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.2.3
- ghcr.io/bnb-chain/op-batcher:v0.2.3
- ghcr.io/bnb-chain/op-proposer:v0.2.3

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.2.2...v0.2.3

## v0.2.2

This is the Fermat Hardfork release for opBNB Mainnet.
It will be activated at block height 9397477, expected to occur on November 28, 2023, at 6 AM UTC.

All mainnet nodes must upgrade to this release before the hardfork.

### User Facing Changes

- Support reading private keys from AWS Secret Manager for `op-node`, `op-batcher`, and `op-proposer`. Refer to PR #72 for additional information.

### Partial Changelog

- #72: feat: support AWS key manager

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.2.2
- ghcr.io/bnb-chain/op-batcher:v0.2.2
- ghcr.io/bnb-chain/op-proposer:v0.2.2

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.2.1...v0.2.2

## v0.2.1

This is a minor release and upgrading is optional.

### User Facing Changes

- Adds a `--rpc.admin-state` CLI option to specify a file to persist config changes made via the RPC Admin APIs to.
- Add `admin_sequencerActive` RPC method. Returns true if the node is actively sequencing, otherwise false.

These features are merged from upstream code. Check the following PRs for more details:
- https://github.com/ethereum-optimism/optimism/pull/6190
- https://github.com/ethereum-optimism/optimism/pull/6105

### Partial Changelog

- #78: feat: support persist active api

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.2.1
- ghcr.io/bnb-chain/op-batcher:v0.2.1
- ghcr.io/bnb-chain/op-proposer:v0.2.1

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.2.0...v0.2.1


## v0.2.0

This is a hardfork release for the opBNB Testnet called Fermat.
It will be activated at block height 12113000, expected to occur on November 3, 2023, at 6 AM UTC.

### User Facing Changes

- The L1 gas price of all L2 transactions will be fixed(3 Gwei by default and 5 Gwei for Testnet).
- Introduce a new type of RPC kind called `bsc_fullnode`. To enable it, include the parameter `--l1.rpckind=bsc_fullnode` if the layer 1 endpoint supports the `eth_getTransactionReceiptsByBlockNumber` API. This will significantly enhance the performance of retrieving L1 receipts (#63).
- The rollup configuration for opBNB Mainnet and Testnet has been added to the code. You can now use the `--network=opBNBTestnet` or `--network=opBNBMainnet` flag instead of `--rollup.config=./rollup.json` to specify the rollup configuration for the op-node. (#65)
-  Allow the addition of multiple L1 endpoints in the configuration. For example: `--l1=https://data-seed-prebsc-1-s1.binance.org:8545,https://data-seed-prebsc-2-s2.binance.org:8545,https://data-seed-prebsc-2-s3.binance.org:8545`. By default, it will use the first endpoint, and if it's unavailable, it will automatically switch to the next one (#55).
- Enable the layer 2 sync mechanism for opBNB by adding `--l2.engine-sync=true` flag on the op-node. Additionally, a new flag ﻿l2.skip-sync-start-check is introduced to allow users to skip the sanity check of L1 origins for unsafe L2 blocks when determining the sync-starting point. (#62)

### Partial Changelog

- #55: feat(op-node/op-batcher/op-proposer): add fallbackClient
- #57: feat(op-node): add pre fetch receipts logic
- #62: feat: engine p2p sync feature
- #63: feat(op-node): add rpcKind bsc_fullnode for eth_getTransactionReceiptsByBlockNumber
- #64: fix(op-e2e): fallback to not use bsc specific method eth_getFinalizedBlock
- #65: feat(op-node): update l1 gas price with fixed value
- #66: feature(contracts-bedrock): add verify.ts for verifying any contract
- #67: fix(op-node): eth_client replace eth_getFinalizedBlock to eth_getFinalizedHeader

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.2.0
- ghcr.io/bnb-chain/op-batcher:v0.2.0
- ghcr.io/bnb-chain/op-proposer:v0.2.0

### Full Changelog

https://github.com/bnb-chain/opbnb/compare/v0.1.2...v0.2.0

## v0.1.2

This is the initial release for opBNB Testnet.

The repo base is [optimism v1.1.0](https://github.com/ethereum-optimism/optimism/releases/tag/op-node%2Fv1.1.0).

### Changelog

1. [feat: support non-eip-1559 l1 & compatible for BSC](https://github.com/bnb-chain/opbnb/commit/2867cfac0a3b4a505e2cac73f7659b0bef5743e5)
2. [feat: op-proposer propose currentL1Hash behavior config by AllowNonFinalizedFlag](https://github.com/bnb-chain/opbnb/commit/19602ccb037073301296875e3c4d4d9d97b8e99c)
3. [feat: ResourceMetering.sol compatible with BSC](https://github.com/bnb-chain/opbnb/commit/2ce30b27b6c2352d330522b8397ed8f8ef72f1a8)
4. [fix(op-batcher): solve race condition of BatchSubmitter](https://github.com/bnb-chain/opbnb/pull/5)
