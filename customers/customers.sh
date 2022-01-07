#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0
TMPFILE=`mktemp`
shopt -s extglob

function _exit(){
    printf "Exiting:%s\n" "$1"
    exit -1
}


: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}

# Where am I?
DIR=${PWD}

setGlobalsForCustomer() {
    # Locate the test-network
    cd "${DIR}/../../test-network"
    env | sort > $TMPFILE

    OVERRIDE_ORG="1"
    . ./scripts/envVar.sh

    parsePeerConnectionParameters 1 2

    # set the fabric config path
    export FABRIC_CFG_PATH="${DIR}/../../config"
    export PATH="${DIR}/../../bin:${PWD}:$PATH"

    env | sort | comm -1 -3 $TMPFILE - | sed -E 's/(.*)=(.*)/export \1="\2"/'

    rm $TMPFILE

    cd "${DIR}"
}

installPackageChaincodeCustomer() {
    rm -rf ms-chaincode.tar.gz
    setGlobalsForCustomer
    peer lifecycle chaincode package ms-chaincode.tar.gz --lang golang --path ./chaincode --label ms_0
    peer lifecycle chaincode install ms-chaincode.tar.gz
    echo "===================== Chaincode is packaged on Customer ===================== "
}

queryInstalled() {
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME_1}_${VERSION_1}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo "===================== Query installed successful on Customer on channel ===================== "
}

approveForMyOrg() {
    peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name medicinecontract -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA
    echo "===================== Chaincode approved from org 1 ===================== "
}

installPackageChaincodeCustomer
queryInstalled
approveForMyOrg