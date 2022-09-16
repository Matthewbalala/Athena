#!/bin/bash
source base_parameters.sh

export FABRIC_LOGGING_SPEC=INFO
export ORDERER_GENERAL_LOCALMSPDIR=${FABRIC_CFG_PATH}/crypto-config/ordererOrganizations/${ORDERER_DOMAIN}/orderers/${ORDERER_ADDRESS}.${ORDERER_DOMAIN}/msp

(cd ${FABRIC_ROOT} && make orderer)
orderer start
