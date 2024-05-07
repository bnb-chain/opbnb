#!/usr/bin/env bash
set -eou
cd /db/node-deploy

echo "starting..."
bash -x ./start_cluster.sh start

while true; do
    sleep 1000
done
