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

if [ ! -d node-deploy ]; then
  git clone https://github.com/bnb-chain/node-deploy.git
else
  echo "node-deploy exists"
fi
cd node-deploy
git submodule update --init --recursive
cd genesis
npm install
cd ..

if [ ! -d bsc ]; then
  git clone https://github.com/bnb-chain/bsc.git
else
  echo "bsc exists"
fi

if [ ! -f bin/geth ]; then
  cd bsc && make geth
  cp ./build/bin/geth ../bin/geth
  if [ ! -f ../bin/bootnode ]; then
    go build -o ./build/bin/bootnode ./cmd/bootnode
    cp ./build/bin/bootnode ../bin/bootnode
  else
    echo "bin/bootnode exists"
  fi
  cd ..
else
  echo "bin/geth exists"
fi

if [ ! -d node ]; then
  git clone https://github.com/bnb-chain/node.git
else
  echo "node exists"
fi

if [ ! -f bin/tbnbcli ]; then
  export CGO_CFLAGS="-O -D__BLST_PORTABLE__"
  export CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"
  cd node && make build
  cp ./build/tbnbcli ../bin/tbnbcli
  cp ./build/bnbchaind ../bin/bnbchaind
  cd ..
else
  echo "bin/tbnbcli exists"
fi

if [ ! -d test-crosschain-transfer ]; then
  git clone https://github.com/bnb-chain/test-crosschain-transfer.git
else
  echo "test-crosschain-transfer exists"
fi

if [ ! -f bin/test-crosschain-transfer ]; then
  cd test-crosschain-transfer && go build
  cp ./test-crosschain-transfer ../bin/test-crosschain-transfer
  cd ..
else
  echo "bin/test-crosschain-transfer exists"
fi

if [ ! -f bin/tool ]; then
  make tool
else
  echo "bin/tool exists"
fi

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
