# bedrock-devnet

This is a utility for running a local Bedrock devnet.

It allows us to quickly start the devnet locally (with L1 network as BSC network and L2 network as opbnb network).

# requirement

docker, nodejs 16+, yarn, foundry, python3, pnpm, poetry, go, jq

Tips:

Install Foundry by following [the instructions located here](https://getfoundry.sh/).

Please make sure your Foundry version matches the one described in versions.json.
If they do not match, please use a command such as `foundryup -C xxxxxx` to modify it.

# usage
First, execute `pnpm install` and `pnpm build` commands in the root directory.

Then we can use the following commands in the project root directory:

Initialize and start devnet:

```
make devnet-up

```

Stop devnet:

```
make devnet-down

```

Stop and clean devnet data:

```
make devnet-clean

```

View devnet logs:

```
make devnet-logs

```

# Notes
1. When executing for the first time, please be patient if you see the message "Waiting for RPC server at...", as the BSC network takes time to initialize.
2. If you encounter an error during the "Deploying contracts" step, please try again as it usually recovers.
3. L1 is accessible at http://localhost:8545, and L2 is accessible at http://localhost:9545

# Additional Information
L1 chain ID is 714.

L2 chain ID is 901.

L1 test account:

address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

Private key: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

L2 test account:

Address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266

Private key: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

The easiest way to invoke this script is to run `make devnet-up` from the root of this repository. Otherwise, to use this script run `python3 main.py --monorepo-dir=<path to the monorepo>`. You may need to set `PYTHONPATH` to this directory if you are invoking the script from somewhere other than `bedrock-devnet`.
