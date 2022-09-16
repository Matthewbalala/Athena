#!/usr/bin/env bash
source base_parameters.sh

export CORE_PEER_ADDRESS=${FAST_PEER_ADDRESS}:7051
peer chaincode query -C ${CHANNEL} -n ${CHAINCODE} -c '{"Args":["query","'${1}'"]}'