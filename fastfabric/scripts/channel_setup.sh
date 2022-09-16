#!/bin/bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp
export CORE_PEER_ADDRESS=${FAST_PEER_ADDRESS}:7051

peer channel create -o ${ORDERER_ADDRESS}:7050 -c ${CHANNEL} -f ./channel-artifacts/channel.tx

if [[ ! -z $STORAGE_ADDRESS ]];then
    export CORE_PEER_ADDRESS=$(get_correct_peer_address ${STORAGE_ADDRESS}):7051
    peer channel join -b ${CHANNEL}.block
fi

for i in ${ENDORSER_ADDRESS[@]}
do
    export CORE_PEER_ADDRESS=$(get_correct_peer_address ${i}):7051
    peer channel join -b ${CHANNEL}.block
done

export CORE_PEER_ADDRESS=$(get_correct_peer_address ${FAST_PEER_ADDRESS}):7051
peer channel join -b ${CHANNEL}.block
peer channel update -o $(get_correct_orderer_address):7050 -c ${CHANNEL} -f ./channel-artifacts/anchor_peer.tx
