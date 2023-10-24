# bedrock-devnet

This is a utility for running a local Bedrock devnet.

It allows us to quickly start the devnet locally (with L1 network as BSC network and L2 network as opbnb network).

# requirement

docker, nodejs 16+, yarn, foundry, python2, python3

Tips:

Install Foundry by following [the instructions located here](https://getfoundry.sh/).

# usage
First, execute `yarn install` and `yarn build` commands in the root directory.

Then we can use the following commands in the project root directory:

Initialize and start devnet:

```
make devnet-up-deploy

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
1. If you encounter a ValueError: invalid mode: 'rU' while trying to load binding.gyp error when executing `yarn install`, this may be caused by python3 installed on your computer. You need to install python 2.7 and configure the environment variable to specify the python version to use: `export npm_config_python=/path/to/executable/python`.
2. When executing for the first time, please be patient if you see the message "wait L1 up...", as the BSC network takes time to initialize.
3. If you encounter an error during the "Deploying contracts" step, please try again as it usually recovers.
4. Do not use the `make devnet-up` command, use the make `devnet-up-deploy` command to start devnet. The `devnet-up` command is not well adapted.
5. L1 is accessible at http://localhost:8545, and L2 is accessible at http://localhost:9545

# Additional Information
L1 chain ID is 714.

L2 chain ID is 901.

L1 test account:

address: 0x04d63aBCd2b9b1baa327f2Dda0f873F197ccd186

Private key: 59ba8068eb256d520179e903f43dacf6d8d57d72bd306e1bd603fdb8c8da10e8

L2 test account:

Address: 0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266

Private key: ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

