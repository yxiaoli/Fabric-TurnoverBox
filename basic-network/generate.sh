#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
#export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
#export FABRIC_CFG_PATH=${PWD}
# export PATH=/home/sherry/Code/fabric-samples/bin:$PATH
CHANNEL_NAME=mychannel

docker-compose -f docker-compose-cli.yaml down --volumes --remove-orphans
docker rm -f $(docker ps -aq)
# remove previous crypto material and config transactions
rm -fr channel-artifacts/*
rm -fr crypto-config/*

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# generate genesis block for orderer
#export PATH=/home/sherry/Code/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=${PWD}
configtxgen -profile FourOrgOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

# generate channel configuration transaction
export CHANNEL_NAME=mychannel
configtxgen -profile FourOrgChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel configuration transaction..."
  exit 1
fi

# generate anchor peer transaction
configtxgen -profile FourOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/OperatorMSPanchors.tx -channelID $CHANNEL_NAME -asOrg OperatorMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for OperatorMSP..."
  exit 1
fi

configtxgen -profile FourOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/SupplierMSPanchors.tx -channelID $CHANNEL_NAME -asOrg SupplierMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for SupplierMSP..."
  exit 1
fi

configtxgen -profile FourOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/DistributorMSPanchors.tx -channelID $CHANNEL_NAME -asOrg DistributorMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for DistributorMSP..."
  exit 1
fi

configtxgen -profile FourOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/RetailerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg RetailerMSP
if [ "$?" -ne 0 ]; then
  echo "Failed to generate anchor peer update for RetailerMSP..."
  exit 1
fi