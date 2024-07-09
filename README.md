# opBNB - High-performance layer 2 solution

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is Optimism?](#what-is-optimism)
- [Documentation](#documentation)
- [Specification](#specification)
- [Community](#community)
- [Contributing](#contributing)
- [Security Policy and Vulnerability Reporting](#security-policy-and-vulnerability-reporting)
- [Directory Structure](#directory-structure)
- [Development and Release Process](#development-and-release-process)
  - [Overview](#overview)
  - [Production Releases](#production-releases)
  - [Development branch](#development-branch)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is Optimism?

It works by offloading transaction processing and resource usage from the BNB Smart Chain, while still posting data to the underlying network. Users interact with the opBNB network by depositing funds from BSC and using applications and contracts on opBNB. At its core, opBNB allows users to deposit and withdraw funds, use smart contracts, and view network data with high throughput and low fees. By leveraging Layer 2, opBNB is able to scale beyond the constraints of the BNB Smart Chain and provide an improved experience for users.

## Comparison

Besides the [differentiators of bedrock](https://community.optimism.io/docs/developers/bedrock/differences/), opBNB is the solution that we aim to provide the best optimistic rollup solution on the BSC.

- Capacity can reach to 100m gas per second, which is much higher than other layer 2 solutions on the Ethereum.
- Gas fee of transfer can reach as low as $0.005 on average.
- block time is 1 second.

| **Parameter**                         | **opBNB value** | **Optimism value** | **Ethereum value (for reference)** |
| ------------------------------------- | --------------- | ------------------ | ---------------------------------- |
| Block gas limit                       | 100,000,000 gas | 30,000,000 gas     | 30,000,000 gas                     |
| Block gas target                      | 50,000,000      | 5,000,000 gas      | 15,000,000 gas                     |
| EIP-1559 elasticity multiplier        | 2               | 6                  | 2                                  |
| EIP-1559 denominator                  | 8               | 50                 | 8                                  |
| Maximum base fee increase (per block) | 12.5%           | 10%                | 12.5%                              |
| Maximum base fee decrease (per block) | 12.5%           | 2%                 | 12.5%                              |
| Block time in seconds                 | 1               | 2                  | 12                                 |

## Documentation

- If you want to build on top of OP Mainnet, refer to the [Optimism Documentation](https://docs.optimism.io)
- If you want to build your own OP Stack based blockchain, refer to the [OP Stack Guide](https://docs.optimism.io/stack/getting-started), and make sure to understand this repository's [Development and Release Process](#development-and-release-process)

## Specification

If you're interested in the technical details of how Optimism works, refer to the [Optimism Protocol Specification](https://github.com/ethereum-optimism/specs).

## Community

To get help from other developers, discuss ideas, and stay up-to-date on what's happening, become a part of our community on Discord. Join our [official Discord Channel](https://discord.com/invite/bnbchain).

## Contributing

Read through [CONTRIBUTING.md](./CONTRIBUTING.md) for a general overview of the contributing process for this repository.
Use the [Developer Quick Start](./CONTRIBUTING.md#development-quick-start) to get your development environment set up to start working on the Optimism Monorepo.
Then check out the list of [Good First Issues](https://github.com/ethereum-optimism/optimism/issues?q=is:open+is:issue+label:D-good-first-issue) to find something fun to work on!
Typo fixes are welcome; however, please create a single commit with all of the typo fixes & batch as many fixes together in a PR as possible. Spammy PRs will be closed.

## Security Policy and Vulnerability Reporting

Please refer to the canonical [Security Policy](https://github.com/ethereum-optimism/.github/blob/master/SECURITY.md) document for detailed information about how to report vulnerabilities in this codebase.
Bounty hunters are encouraged to check out [the Optimism Immunefi bug bounty program](https://immunefi.com/bounty/optimism/).
The Optimism Immunefi program offers up to $2,000,042 for in-scope critical vulnerabilities.

## Directory Structure

<pre>
├── <a href="./docs">docs</a>: A collection of documents including audits and post-mortems
├── <a href="./op-batcher">op-batcher</a>: L2-Batch Submitter, submits bundles of batches to L1
├── <a href="./op-e2e">op-e2e</a>: End-to-End testing of all bedrock components in Go
├── <a href="./op-heartbeat">op-heartbeat</a>: Heartbeat monitor service
├── <a href="./op-node">op-node</a>: rollup consensus-layer client
├── <a href="./op-preimage">op-preimage</a>: Go bindings for Preimage Oracle
├── <a href="./op-program">op-program</a>: Fault proof program
├── <a href="./op-proposer">op-proposer</a>: L2-Output Submitter, submits proposals to L1
├── <a href="./op-service">op-service</a>: Common codebase utilities
├── <a href="./op-ufm">op-ufm</a>: Simulations for monitoring end-to-end transaction latency
├── <a href="./op-wheel">op-wheel</a>: Database utilities
├── <a href="./ops">ops</a>: Various operational packages
├── <a href="./ops-bedrock">ops-bedrock</a>: Bedrock devnet work
├── <a href="./packages">packages</a>
│   ├── <a href="./packages/chain-mon">chain-mon</a>: Chain monitoring services
│   ├── <a href="./packages/contracts-bedrock">contracts-bedrock</a>: Bedrock smart contracts
│   ├── <a href="./packages/sdk">sdk</a>: provides a set of tools for interacting with Optimism
├── <a href="./proxyd">proxyd</a>: Configurable RPC request router and proxy
├── <a href="./specs">specs</a>: Specs of the rollup starting at the Bedrock upgrade
</pre>

## Development and Release Process

### Overview

Please read this section if you're planning to fork this repository, or make frequent PRs into this repository.

### Production Releases

Production releases are always tags, versioned as `<component-name>/v<semver>`.
For example, an `op-node` release might be versioned as `op-node/v1.1.2`, and  smart contract releases might be versioned as `op-contracts/v1.0.0`.
Release candidates are versioned in the format `op-node/v1.1.2-rc.1`.
We always start with `rc.1` rather than `rc`.

For contract releases, refer to the GitHub release notes for a given release, which will list the specific contracts being released—not all contracts are considered production ready within a release, and many are under active development.

Tags of the form `v<semver>`, such as `v1.1.4`, indicate releases of all Go code only, and **DO NOT** include smart contracts.
This naming scheme is required by Golang.
In the above list, this means these `v<semver` releases contain all `op-*` components, and exclude all `contracts-*` components.

`op-geth` embeds upstream geth’s version inside it’s own version as follows: `vMAJOR.GETH_MAJOR GETH_MINOR GETH_PATCH.PATCH`.
Basically, geth’s version is our minor version.
For example if geth is at `v1.12.0`, the corresponding op-geth version would be `v1.101200.0`.
Note that we pad out to three characters for the geth minor version and two characters for the geth patch version.
Since we cannot left-pad with zeroes, the geth major version is not padded.

See the [Node Software Releases](https://docs.optimism.io/builders/node-operators/releases) page of the documentation for more information about releases for the latest node components.
The full set of components that have releases are:

- `chain-mon`
- `ci-builder`
- `ci-builder`
- `indexer`
- `op-batcher`
- `op-contracts`
- `op-challenger`
- `op-heartbeat`
- `op-node`
- `op-proposer`
- `op-ufm`
- `proxyd`

All other components and packages should be considered development components only and do not have releases.

### Development branch

The primary development branch is [`develop`](https://github.com/ethereum-optimism/optimism/tree/develop/).
`develop` contains the most up-to-date software that remains backwards compatible with the latest experimental [network deployments](https://community.optimism.io/docs/useful-tools/networks/).
If you're making a backwards compatible change, please direct your pull request towards `develop`.

**Changes to contracts within `packages/contracts-bedrock/src` are usually NOT considered backwards compatible.**
Some exceptions to this rule exist for cases in which we absolutely must deploy some new contract after a tag has already been fully deployed.
If you're changing or adding a contract and you're unsure about which branch to make a PR into, default to using a feature branch.
Feature branches are typically used when there are conflicts between 2 projects touching the same code, to avoid conflicts from merging both into `develop`.

## License

All files within this repository are licensed under the [MIT License](https://github.com/bnb-chain/opbnb/blob/master/LICENSE) unless stated otherwise.
