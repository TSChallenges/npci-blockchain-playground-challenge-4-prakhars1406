#!/bin/bash

# Set Fabric binaries path
export FABRIC_CFG_PATH=$PWD/config/
export PATH=$PWD/bin:$PATH

# Start the network and create channel
function up() {
  echo "Starting network and creating channel"
  docker-compose -f docker-compose-test-net.yaml up -d
  docker exec -it cli bash -c 'peer channel create -o orderer.example.com:7050 -c assetchannel -f /etc/hyperledger/fabric/channel-artifacts/channel.tx'
}

# Deploy Chaincode to the network
function deployCC() {
  echo "Deploying chaincode..."
  docker exec -it cli bash -c 'peer chaincode install -n assetcc -v 1.0 -p /opt/gopath/src/github.com/assetcc'
  docker exec -it cli bash -c 'peer chaincode instantiate -o orderer.example.com:7050 -C assetchannel -n assetcc -v 1.0 -c \'{"Args":[""]}\''
}

case "$1" in
  up)
    up
    ;;
  createChannel)
    up
    ;;
  deployCC)
    deployCC
    ;;
  *)
    echo "Usage: $0 {up|createChannel|deployCC}"
    exit 1
esac