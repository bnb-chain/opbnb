# Changelog

## v0.5.1

This is a minor release and upgrading is optional.

### What's Changed
* fix(ci): support building arm64 architecture by @welkin22 in https://github.com/bnb-chain/opbnb/pull/239
* fix(op-node): l1 client chan stuck when closed in ELSync mode by @welkin22 in https://github.com/bnb-chain/opbnb/pull/241
* fix: add sync status when newpayload API meet the specific error by @krish-nr in https://github.com/bnb-chain/opbnb/pull/240

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.5.1
- ghcr.io/bnb-chain/op-batcher:v0.5.1
- ghcr.io/bnb-chain/op-proposer:v0.5.1

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.5.0...v0.5.1

## v0.5.0

This release includes code merging from the upstream version v1.7.7 along with several fixs and improvements.

Fjord fork from upstream is included. It is set to be activated on both the opBNB Mainnet and Testnet environments according to the following schedule:

- Testnet: Sep-10-2024 06:00 AM +UTC
- Mainnet: Sep-24-2024 06:00 AM +UTC

All mainnet and testnet nodes must upgrade to this release before the hardfork time.
Also note that the `op-geth` should be upgraded to v0.5.0 accordingly, check [this](https://github.com/bnb-chain/op-geth/releases/tag/v0.5.0) for more details.

### User Facing Changes

* The L1 fee calculation is optimized. Check this [spec](https://specs.optimism.io/protocol/fjord/exec-engine.html) for more details.
* New flag `--wait-node-sync` added to op-batcher (default false), indicates if during startup, the batcher should wait for a recent batcher tx on L1 to finalize (via more block confirmations). This should help avoid duplicate batcher txs
* New flag `--wait-node-sync` added to op-proposer (default false), indicates if during startup, the proposer should wait for the rollup node to sync to the current L1 tip before proceeding with its driver loop
* New flag `--compression-algo` added to op-batcher (default zlib), user can choose brotli algo after Fjord fork
* New flag `--l1.rpc-max-cache-size` added to op-node (default 1000), so user can config the the maximum cache size of the L1 client

### What's Changed
* Merge upstream v1.7.7 by @bnoieh in https://github.com/bnb-chain/opbnb/pull/216
* feat(op-node): Keep consistent status when meet an unexpected el sync by @krish-nr in https://github.com/bnb-chain/opbnb/pull/222
* feat(op-node): add l1 cache size config by @welkin22 in https://github.com/bnb-chain/opbnb/pull/225
* feat(op-chain-ops): add Wright fork config into genesis file generation code by @welkin22 in https://github.com/bnb-chain/opbnb/pull/226

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.5.0
- ghcr.io/bnb-chain/op-batcher:v0.5.0
- ghcr.io/bnb-chain/op-proposer:v0.5.0

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.4.4...v0.5.0

## v0.4.4

This release includes important fixes to help pbss geth nodes automatically recover from ungraceful shutdowns.
We recommend upgrading to this version if you are using the pbss mode.

### User Facing Changes

### What's Changed
* fix(devnet): Modify the blob configuration of the bsc devnet by @welkin22 in #212
* feat: auto recover from pbss geth unclean shutdown by @krish-nr in #214
* fix: fix el bug by @krish-nr in #215

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.4.4
- ghcr.io/bnb-chain/op-batcher:v0.4.4
- ghcr.io/bnb-chain/op-proposer:v0.4.4

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.4.3...v0.4.4

## v0.4.3

This is a minor release and upgrading is optional.

### User Facing Changes

* Support choose economic DA type automatically for op-batcher. #209
* Add 2 configs for el-sync optimization and enable fastnode mode again. #201

### What's Changed
* feature: add haber fork config in deployment script by @redhdx in https://github.com/bnb-chain/opbnb/pull/202
* feat(op-node): support multi clients to fetch blobs by @bnoieh in https://github.com/bnb-chain/opbnb/pull/199
* feat: fastnode support by trigger el-sync when needed  by @krish-nr in https://github.com/bnb-chain/opbnb/pull/201
* fix(blob-client): don't append L1ArchiveBlobRpcAddr flag to config if not set by @bnoieh in https://github.com/bnb-chain/opbnb/pull/207
* fix(devnet): fork offset should be 0x by @welkin22 in https://github.com/bnb-chain/opbnb/pull/210
* fix(devnet): batcher uses its address to submit transactions by @welkin22 in https://github.com/bnb-chain/opbnb/pull/211
* feat:  op-batcher auto switch to economic DA type by @bnoieh in https://github.com/bnb-chain/opbnb/pull/209

### Docker Images

- ghcr.io/bnb-chain/op-node:v0.4.3
- ghcr.io/bnb-chain/op-batcher:v0.4.3
- ghcr.io/bnb-chain/op-proposer:v0.4.3

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.4.2...v0.4.3

## v0.4.2

This is the mainnet hardfork release version.

Four hard forks are scheduled to launch on the opBNB Mainnet:
Shanghai/Canyon Time: 2024-06-20 08:00:00 AM UTC
Delta Time: 2024-06-20 08:10:00 AM UTC
Cancun/Ecotone Time: 2024-06-20 08:20:00 AM UTC
Haber Time: 2024-06-20 08:30:00 AM UTC

All mainnet `op-node` have to be upgraded to this version before 2024-06-20 08:00:00 AM UTC.
The `op-geth` also have to be upgraded to v0.4.2 accordingly, check [this](https://github.com/bnb-chain/op-geth/releases/tag/v0.4.2) for more details.

### User Facing Changes

If you are upgrading from v0.3.x to this version, please note that there are some configuration changes.
-  Removed `--l1.rpckind=bsc_fullnode`
-  Removed `--l2.engine-sync`
-  Removed `--l2.skip-sync-start-check`
-  To start engine-sync, use `--syncmode=execution-layer` (default value is `consensus-layer`)
-  Added `--l1.max-concurrency=20` to control the rate of requests to L1 endpoints.

After the Cancun/Ecotone hard fork, DA data will be submitted to the BSC network in blob format. Regular BSC nodes only retain blob data from the past 18 days. If you are syncing data from the genesis block or are more than 18 days behind the latest block, you will need to ensure that your configured L1 endpoint supports persisting blob data for a longer period of time. We will ensure that the snapshot provided by this [snapshot repository](https://github.com/bnb-chain/opbnb-snapshot) is within the 18-day range, so you can also choose to use the snapshot to avoid relying on older blob data to start your new node.

### What's Changed
* feature: update deployment script for opBNB by @redhdx in https://github.com/bnb-chain/opbnb/pull/196
* fix: fix CI after 4844 merge by @welkin22 in https://github.com/bnb-chain/opbnb/pull/198
* op-node: set finalityDelay to 15 to speed up finality update by @bnoieh in https://github.com/bnb-chain/opbnb/pull/200
* config: Mainnet canyon/delta/ecotone fork time by @welkin22 in https://github.com/bnb-chain/opbnb/pull/203

### Docker Images
- ghcr.io/bnb-chain/op-node:v0.4.2
- ghcr.io/bnb-chain/op-batcher:v0.4.2
- ghcr.io/bnb-chain/op-proposer:v0.4.2

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.4.1...v0.4.2

## v0.4.1

This is a minor release and upgrading is optional.

### User Facing Changes

- Add flag `--txmgr.blob-gas-price-limit` for op-batcher to limit the maximum gas price of submitted tx

### Partial Changelog

* fix: fix devnet after 1.7.2 upstream merge by @welkin22 in https://github.com/bnb-chain/opbnb/pull/194
* op-batcher: optimize tx submitting and add metrics by @bnoieh in https://github.com/bnb-chain/opbnb/pull/195

### Docker Images

- ghcr.io/bnb-chain/op-batcher:v0.4.1

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.4.0...v0.4.1

## v0.4.0

This release includes code merging from the upstream version v1.7.2 to transition Testnet's DA data from calldata to blob format.

Four hard forks are scheduled to launch on the opBNB Testnet:
Snow Time: May-15-2024 06:00 AM +UTC
Shanghai/Canyon Time: May-15-2024 06:10 AM +UTC
Delta Time: May-15-2024 06:20 AM +UTC
Cancun/Ecotone Time: May-15-2024 06:30 AM +UTC

### User Facing Changes
Nodes on the **Testnet** need to be upgraded to this version before the first hard fork time.
**Note: This is a version prepared for Testnet, Mainnet nodes do not need to upgrade to this version.**

**Note: After the Cancun/Ecotone hard fork, DA data will be submitted to the BSC network in blob format. Regular BSC nodes only retain blob data from the past 18 days. If you are syncing data from the genesis block or are more than 18 days behind the latest block, you will need to ensure that your configured L1 endpoint supports persisting blob data for a longer period of time. We will ensure that the Testnet snapshot provided by this [snapshot repository](https://github.com/bnb-chain/opbnb-snapshot) is within the 18-day range, so you can also choose to use the snapshot to avoid relying on older blob data to start your new node.**

Changes in op-node configuration:
-  Removed `--l1.rpckind=bsc_fullnode`
-  Removed `--l2.engine-sync`
-  Removed `--l2.skip-sync-start-check`
-  To start engine-sync, use `--syncmode=execution-layer` (default value is `consensus-layer`)
-  Added `--l1.max-concurrency=20` to control the rate of requests to L1 endpoints.

### What's Changed
* feature(op-node): update opBNB qanet info by @redhdx in https://github.com/bnb-chain/opbnb/pull/187
* feat: update qanet config by @redhdx in https://github.com/bnb-chain/opbnb/pull/188
* feature(op-node): add opBNB qanet hard fork config by @redhdx in https://github.com/bnb-chain/opbnb/pull/189
* Fix blob parsing problem by @welkin22 in https://github.com/bnb-chain/opbnb/pull/190
* chore: fork config for 4844-2 qanet by @welkin22 in https://github.com/bnb-chain/opbnb/pull/191
* Merge upstream v1.7.2 by @bnoieh in https://github.com/bnb-chain/opbnb/pull/184
* config: Testnet 4844 fork time by @welkin22 in https://github.com/bnb-chain/opbnb/pull/192

### Docker Images
ghcr.io/bnb-chain/op-node:v0.4.0
ghcr.io/bnb-chain/op-batcher:v0.4.0
ghcr.io/bnb-chain/op-proposer:v0.4.0

**Full Changelog**: https://github.com/bnb-chain/opbnb/compare/v0.3.3...v0.4.0

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
