#!/usr/bin/env bash
set -eou
cd /
export workspace=/db
export genesisDir=/genesis
apk add expect wget nodejs npm git
if [ ! -f solc-linux-amd64-v0.6.4+commit.1dca32f3 ]; then
  wget https://github.com/ethereum/solc-bin/raw/gh-pages/linux-amd64/solc-linux-amd64-v0.6.4%2Bcommit.1dca32f3
  cp solc-linux-amd64-v0.6.4+commit.1dca32f3 /usr/bin/solc
  chmod +x /usr/bin/solc
else
  echo "solc already exists"
fi

if [ ! -d ${workspace}/bsc ]; then
  mkdir -p ${workspace}/bsc/validator0
  echo "${KEYPASS}" > ${workspace}/bsc/password.txt
  cons_addr=$(geth account new --datadir ${workspace}/bsc/validator0 --password ${workspace}/bsc/password.txt | grep "Public address of the key:" | awk -F"   " '{print $2}')
  fee_addr=$(geth account new --datadir ${workspace}/bsc/validator0_fee --password ${workspace}/bsc/password.txt | grep "Public address of the key:" | awk -F"   " '{print $2}')
  mkdir -p ${workspace}/bsc/bls0
  expect create_bls_key.sh ${workspace}/bsc/bls0
else
   echo "${workspace}/bsc exists."
fi

function prepare_config() {
    git clone https://github.com/bnb-chain/node-deploy.git \
        && cd node-deploy && git submodule update --init --recursive
    cp -R genesis ${genesisDir}
    cd ${genesisDir}
    npm install
    rm -f ${genesisDir}/validators.conf
    rm -f ${genesisDir}/init_holders.template
    cp /init_holders.template ${genesisDir}/init_holders.template

    sed -i -e "s/${replaceWhitelabelRelayer}/${INIT_HOLDER}/g" ${genesisDir}/contracts/RelayerHub.template
    sed -i -e "s/function whitelistInit() external/function whitelistInit() public/g" ${genesisDir}/contracts/RelayerHub.template
    sed -i -e "s/alreadyInit = true;/whitelistInit();\nalreadyInit = true;/g" ${genesisDir}/contracts/RelayerHub.template
    sed -i -e "s/alreadyInit = true;/enableMaliciousVoteSlash = true;\nalreadyInit = true;/g" ${genesisDir}/contracts/SlashIndicator.template
    sed -i -e "s/numOperator = 2;/operators[VALIDATOR_CONTRACT_ADDR] = true;\noperators[SLASH_CONTRACT_ADDR] = true;\nnumOperator = 4;/g" ${genesisDir}/contracts/SystemReward.template
    sed -i -e "s/for (uint i; i<validatorSetPkg.validatorSet.length; ++i) {/ValidatorExtra memory validatorExtra;\nfor (uint i; i<validatorSetPkg.validatorSet.length; ++i) {\n validatorExtraSet.push(validatorExtra);\n validatorExtraSet[i].voteAddress=validatorSetPkg.voteAddrs[i];/g" ${genesisDir}/contracts/BSCValidatorSet.template
    sed -i -e "s/false/true/g" ${genesisDir}/generate-relayerhub.js
    sed -i -e "s/\"0x\" + publicKey.pop()/vs[4]/g" ${genesisDir}/generate-validator.js
    sed "s/{{INIT_HOLDER_ADDR}}/${INIT_HOLDER}/g" ${genesisDir}/init_holders.template > ${genesisDir}/init_holders.js
    for f in ${workspace}/bsc/validator0/keystore/*;do
        cons_addr="0x$(cat ${f} | jq -r .address)"
    done

    for f in ${workspace}/bsc/validator0_fee/keystore/*;do
        fee_addr="0x$(cat ${f} | jq -r .address)"
    done

    mkdir -p ${workspace}/bsc/clusterNetwork/node0
    bbcfee_addrs=${fee_addr}
    powers="0x000001d1a94a2000"
    mv ${workspace}/bsc/bls0/bls ${workspace}/bsc/clusterNetwork/node0/ && rm -rf ${workspace}/bsc/bls0
    vote_addr=0x$(cat ${workspace}/bsc/clusterNetwork/node0/bls/keystore/*json| jq .pubkey | sed 's/"//g')
    echo "${cons_addr},${bbcfee_addrs},${fee_addr},${powers},${vote_addr}" >> ${genesisDir}/validators.conf
    echo "validator" ":" ${cons_addr}
    echo "validatorFee" ":" ${fee_addr}
    #fix l1 genesis timestamp for first block
    timestamp=$(date +"%s")
    hex_timestamp=$(printf "0x%x" $timestamp)
    sed -i -e 's/"timestamp": "0x5e9da7ce"/"timestamp": "'$hex_timestamp'"/g' ./genesis-template.json

    node generate-validator.js
    node generate-genesis.js --chainid ${BSC_CHAIN_ID} --bscChainId "$(printf '%04x\n' ${BSC_CHAIN_ID})"
}

if [ ! -f ${genesisDir}/genesis.json ]; then
  prepare_config
else
   echo "genesis.json exists."
fi

function generate() {
    cd /
    geth init-network --init.dir ${workspace}/bsc/clusterNetwork --init.size=1 --config /config.toml ${genesisDir}/genesis.json
    staticPeers=""
    line=`grep -n -e 'StaticNodes' ${workspace}/bsc/clusterNetwork/node0/config.toml | cut -d : -f 1`
    head -n $((line-1)) ${workspace}/bsc/clusterNetwork/node0/config.toml >> ${workspace}/bsc/clusterNetwork/node0/config.toml-e
    echo "StaticNodes = [${staticPeers}]" >> ${workspace}/bsc/clusterNetwork/node0/config.toml-e
    tail -n +$(($line+1)) ${workspace}/bsc/clusterNetwork/node0/config.toml >> ${workspace}/bsc/clusterNetwork/node0/config.toml-e
    rm -f ${workspace}/bsc/clusterNetwork/node0/config.toml
    mv ${workspace}/bsc/clusterNetwork/node0/config.toml-e ${workspace}/bsc/clusterNetwork/node0/config.toml

    sed -i -e "s/TriesInMemory = 0/TriesInMemory = 128/g" ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e "s/NetworkId = 714/NetworkId = ${BSC_CHAIN_ID}/g" ${workspace}/bsc/clusterNetwork/node0/config.toml

    sed -i -e '/BerlinBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/EWASMBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/CatalystBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/YoloV3Block/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/LondonBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/ArrowGlacierBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/MergeForkBlock/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/TerminalTotalDifficulty/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/BaseFee/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e '/RPCTxFeeCap/d' ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e "s/MirrorSyncBlock = 1/MirrorSyncBlock = 0/g" ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e "s/BrunoBlock = 1/BrunoBlock = 0/g" ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e "s/EulerBlock = 2/EulerBlock = 0\nNanoBlock = 0/g" ${workspace}/bsc/clusterNetwork/node0/config.toml
    sed -i -e 's/DataDir/BLSPasswordFile = \"{{BLSPasswordFile}}\"\nDataDir/g' ${workspace}/bsc/clusterNetwork/node0/config.toml
    PassWordPath="${workspace}/bsc/password.txt"
    sed -i -e "s:{{BLSPasswordFile}}:${PassWordPath}:g" ${workspace}/bsc/clusterNetwork/node0/config.toml
}

if [ ! -d ${workspace}/bsc/clusterNetwork/node0/geth/chaindata ]; then
  generate
else
   echo "geth already init."
fi



export HTTPPort=8545
export WSPort=8545
export MetricsPort=6060
for ((i=0;i<1;i++));do
        cp -R ${workspace}/bsc/validator${i}/keystore ${workspace}/bsc/clusterNetwork/node${i}
        for j in ${workspace}/bsc/validator${i}/keystore/*;do
            cons_addr="0x$(cat ${j} | jq -r .address)"
        done

        # sorry for magic
        for ((k=0;k<1;k++));do
            p2p_port_k=$((30311 + k))
            if [ ${k} -ne ${i} ];then
                sed -i.bak "s/bsc-node-${k}.bsc.svc.cluster.local:30311/localhost:${p2p_port_k}/" ${workspace}/bsc/clusterNetwork/node${i}/config.toml
            else
                sed -i.bak "s/\":30311/\":${p2p_port_k}/" ${workspace}/bsc/clusterNetwork/node${i}/config.toml
            fi
        done
    done

# Start op-geth.
exec geth \
  --config ${workspace}/bsc/clusterNetwork/node0/config.toml \
  --datadir ${workspace}/bsc/clusterNetwork/node0 \
  --password ${workspace}/bsc/password.txt \
  --nodekey ${workspace}/bsc/clusterNetwork/node0/geth/nodekey \
  -unlock ${cons_addr} --rpc.allow-unprotected-txs --allow-insecure-unlock  \
  --ws.addr 0.0.0.0 --ws.port ${WSPort} --http.addr 0.0.0.0 --http.port ${HTTPPort} --http.corsdomain "*" \
  --metrics --metrics.addr localhost --metrics.port ${MetricsPort} --metrics.expensive \
  --gcmode archive --syncmode=full --mine --vote \
