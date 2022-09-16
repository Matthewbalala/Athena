#!/usr/bin/env bash
source base_parameters.sh

if [[ ! -f core.yaml.bak ]]; then
    cp core.yaml core.yaml.bak
fi

bootstrap="bootstrap: $(get_correct_peer_address ${FAST_PEER_ADDRESS}):7051"
(cat core.yaml.bak | sed "s/bootstrap: 127.0.0.1:7051/$bootstrap/g") > core.yaml

if [[ ! -f crypto-config.yaml.bak ]]; then
    cp crypto-config.yaml crypto-config.yaml.bak
fi

endorsers=""
for i in ${ENDORSER_ADDRESS[@]}
do
    endorsers="- Hostname: ${i}\n    ${endorsers}"
done

storage=""
if [[ ! -z $STORAGE_ADDRESS ]]; then
    storage="- Hostname: ${STORAGE_ADDRESS}\n    ${storage}"
fi

(cat crypto-config.yaml.bak | sed "s/ORDERER_DOMAIN/$ORDERER_DOMAIN/g" | sed "s/ORDERER_ADDRESS/$ORDERER_ADDRESS/g"| sed "s/PEER_DOMAIN/$PEER_DOMAIN/g"| sed "s/FAST_PEER_ADDRESS/$FAST_PEER_ADDRESS/g"| sed "s/ENDORSERS/$endorsers/g"| sed "s/STORAGE/$storage/g") > crypto-config.yaml

if [[ ! -f configtx.yaml.bak ]]; then
    cp configtx.yaml configtx.yaml.bak
fi

(cat configtx.yaml.bak | sed "s/ORDERER_DOMAIN/$ORDERER_DOMAIN/g" | sed "s/ORDERER_ADDRESS/$ORDERER_ADDRESS/g"| sed "s/PEER_DOMAIN/$PEER_DOMAIN/g"| sed "s/FAST_PEER_ADDRESS/$FAST_PEER_ADDRESS/g") > configtx.yaml

if [[ -d ./crypto-config ]]; then rm -r ./crypto-config; fi
if [[ -d ./channel-artifacts ]]; then rm -r ./channel-artifacts; fi
mkdir channel-artifacts
$FABRIC_ROOT/.build/bin/cryptogen generate --config=crypto-config.yaml
$FABRIC_ROOT/.build/bin/configtxgen -configPath ./ -outputBlock ./channel-artifacts/genesis.block -profile OneOrgOrdererGenesis -channelID ${CHANNEL}-system-channel
$FABRIC_ROOT/.build/bin/configtxgen -configPath ./ -outputCreateChannelTx ./channel-artifacts/channel.tx -profile OneOrgChannel -channelID ${CHANNEL}
$FABRIC_ROOT/.build/bin/configtxgen -configPath ./ -outputAnchorPeersUpdate ./channel-artifacts/anchor_peer.tx -profile OneOrgChannel -asOrg Org1MSP -channelID ${CHANNEL}
