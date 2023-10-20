#!/usr/bin/env bash
set -eou
cd /db
if [ ! -f solc-linux-amd64-v0.6.4+commit.1dca32f3 ]; then
  wget https://github.com/ethereum/solc-bin/raw/gh-pages/linux-amd64/solc-linux-amd64-v0.6.4%2Bcommit.1dca32f3
  cp solc-linux-amd64-v0.6.4+commit.1dca32f3 /usr/bin/solc
  chmod +x /usr/bin/solc
else
  echo "solc already exists"
fi

cd node-deploy
git submodule update --init --recursive
cd genesis
npm install
cd ..

if [ ! -f init_file_bc ]; then
  bash +x ./setup_bc_node.sh native_init
  echo "finish" > init_file_bc
else
  echo "bc init already finish"
fi
bash +x ./setup_bc_node.sh native_start

if [ ! -f init_file_bsc ]; then
  bash +x ./setup_bsc_node.sh native_init
  echo "finish" > init_file_bsc
else
  echo "bsc init already finish"
fi
bash +x ./setup_bsc_node.sh native_start


while true; do
    sleep 1000
done
