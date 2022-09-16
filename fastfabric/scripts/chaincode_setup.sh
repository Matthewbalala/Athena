#!/bin/bash
source base_parameters.sh

export CORE_PEER_MSPCONFIGPATH=./crypto-config/peerOrganizations/${PEER_DOMAIN}/users/Admin@${PEER_DOMAIN}/msp

endorsers=(${FAST_PEER_ADDRESS})

if [[ ! -z ${ENDORSER_ADDRESS[@]} ]]
then
    endorsers=(${endorsers[@]} ${ENDORSER_ADDRESS[@]})
fi

for i in ${endorsers[@]}
do
    export CORE_PEER_ADDRESS=$(get_correct_peer_address ${i}):7051
    peer chaincode install -l golang -n ${CHAINCODE} -v 1.0 -o ${ORDERER_ADDRESS}:7050 -p "github.com/hyperledger/fabric/fastfabric/chaincode"
done

export CORE_PEER_ADDRESS=$(get_correct_peer_address ${endorsers[0]}):7051
a="'{\"Args\":[\"init\",\"0\", \"1\", \"0\"]}'"
echo peer chaincode instantiate -o $(get_correct_orderer_address):7050 -C ${CHANNEL} -n $CHAINCODE -v 1.0 -c ${a} | bash

sleep 5

bash chaincode_account_setup.sh $1 $2 $3
