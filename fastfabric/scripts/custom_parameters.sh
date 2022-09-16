#!/usr/bin/env bash


export PEER_DOMAIN=""           # can be anything if running on localhost
export ORDERER_DOMAIN=""        # can be anything if running on localhost

# fill in the addresses without domain suffix and without ports
export FAST_PEER_ADDRESS=""
export ORDERER_ADDRESS=""

# leave endorser address and storage address blank if you want to run on a single server
export ENDORSER_ADDRESS=()      # can define multiple addresses in the form ( "addr1" "addr2" ... )
export STORAGE_ADDRESS=""

export CHANNEL=""               # the name of the created channel of the network
export CHAINCODE=""             # the name of the chaincode used on the channel