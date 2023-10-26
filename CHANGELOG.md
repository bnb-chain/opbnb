# Changelog

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
