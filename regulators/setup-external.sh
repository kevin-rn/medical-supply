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

# remove old keystore and wallet inside regulators application
cleanUpCredentials() {
    rm -rf "${DIR}/application/keystore"
    rm -rf "${DIR}/application/wallet"
    rm -rf "${DIR}/application/tpmkey.txt"
}

# Sets certain environment variables in command window (administrator) in order to use the correct set of peer binaries, 
# send commands to the address of the organisation peer, and sign requests with the correct cryptographic material.
setGlobalsForRegulator() {
    # Locate the test-network
    cd "${DIR}/../../test-network"
    env | sort > $TMPFILE

    OVERRIDE_ORG="2"
    . ./scripts/envVar.sh

    parsePeerConnectionParameters 1 2

    # set the fabric config path
    export FABRIC_CFG_PATH="${DIR}/../config"
    export PATH="${DIR}/../../bin:${PWD}:$PATH"

    env | sort | comm -1 -3 $TMPFILE - | sed -E 's/(.*)=(.*)/export \1="\2"/'

    rm $TMPFILE

    cd "${DIR}"
}

# Package the smart contract into a chaincode and installs it.
installPackageChaincodeRegulator() {
    rm -rf ms-chaincode.tar.gz

    # Set global enviroments
    setGlobalsForRegulator

    # Install chaincode
    peer lifecycle chaincode install ../external-chaincode/ms-external-chaincode.tar.gz
    echo "===================== Chaincode is packaged on Regulator ===================== "
}

# Query the installed chaincode to get the package_id and sets it as an environmental variable.
queryInstalled() {
    peer lifecycle chaincode queryinstalled >&log.txt
    cat log.txt
    PACKAGE_ID=$(sed -n "/${CC_NAME_1}_${VERSION_1}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
    echo "===================== Query installed successful on Regulator on channel ===================== "
}

# Approve chaincode for the organisation.
approveForMyOrg() {
    peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name medicinecontract -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA
    echo "===================== Chaincode approved from org 1 ===================== "
}

# Commit the chaincode definition
commitChaincodeDefinition() {
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name medicinecontract -v 0 --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
    echo "===================== Chaincode definition committed ===================== "
}

cleanUpCredentials
installPackageChaincodeRegulator
queryInstalled
approveForMyOrg
commitChaincodeDefinition