# opbnb-upgrade

This document introduces how to use multiple scripts and cli tools to upgrade the opbnb contracts.

### Upgrade Steps

#### Deploy new contracts

```
cd ~/opbnb/packages/contracts-bedrock

export IMPL_SALT=$(openssl rand -hex 32)
export DEPLOY_CONFIG_PATH=~/opbnb/packages/contracts-bedrock/deployconfig/devnetL1-template.json (just for compatible, no use config)
export DEPLOY_PRIVATE_KEY=
export L1_RPC_URL=
export VERIFIER_URL=
export API_KEY=

forge script scripts/DeployUpgrades.s.sol:Deploy --private-key $DEPLOY_PRIVATE_KEY --broadcast --rpc-url $L1_RPC_URL --slow --verify --verifier-url $VERIFIER_URL --etherscan-api-key $API_KEY
```
This will output $ChainId-deploy.json file in opbnb/packages/contracts-bedrock/deployments directory. The file contains new contracts addresses.

#### Check old contracts

```
cd ~/opbnb/op-chain-ops

make opbnb-upgrade

./bin/opbnb-upgrade \
 --l1-rpc-url \
 --old_contracts_check=true \
 --outfile
```

This will output old contracts information in outfile. The chainid will be obtained based on the l1-rpc-url and the corresponding old network contracts information will be output.

#### Generate tx builder json file

```
./bin/opbnb-upgrade \
--l1-rpc-url \
--outfile
```

This will output upgrade tx builder json file.

#### Simulate upgrade

##### generate tx calldata from tx builder json file

```
cd opbnb/packages/contracts-bedrock

// $tx_builder.json is upgrade tx builder which output in previous step. The json file need to be located in fs_permissions path in foundry.toml.
// The sender is one of multisign wallet signers
forge script ./scripts/BuildDataFromJson.s.sol --sig "rawInputData(string memory _path)" $tx_builder.json --via-ir --sender $sender
```

##### forking network

```
anvil --fork-url
```

##### set safe wallet storage

```
anvil --fork-url
```

##### impersonateAccount

```
// params: one of multisign wallet signers address
cast rpc anvil_impersonateAccount $SIGNER
```

##### simulate upgrade tx

```
// params: safe wallet + one of signers + calldata (calldata from tx builder json)
cast send $SAFE_WALLET --unlocked --from $SIGNER $CALLDATA
```

#### Check new contracts

```
cd ~/opbnb/op-chain-ops

make opbnb-upgrade

./bin/opbnb-upgrade \
 --l1-rpc-url \
 --new_contracts_check=true \
 --outfile
```

This will output new contracts information in outfile. The chainid will be obtained based on the l1-rpc-url and the corresponding new network contracts information will be output.

#### Compare contracts

```
cd ~/opbnb/op-chain-ops

make opbnb-upgrade
./bin/opbnb-upgrade \
--l1-rpc-url \
--compare_contracts=true \
--old_contracts_file
--new_contracts_file
```

