#!/usr/bin/env bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp

peer=$1

export CORE_PEER_ADDRESS=${peer}:7051

echo peer chaincode invoke -C ${CHANNEL} -n $CHAINCODE -c '{"Args":["transfer","account1", "account2", "10"]}'

peer chaincode invoke -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["transfer","account1", "account2", "10"]}'

sleep 3

a="'{\"Args\":[\"query\",\"account1\"]}'"
echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash

a="'{\"Args\":[\"query\",\"account2\"]}'"
echo peer chaincode query -C ${CHANNEL}  -n ${CHAINCODE} -c $a | bash
