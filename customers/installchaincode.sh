#!/bin/bash

TMPFILE=`mktemp`
shopt -s extglob

function _exit(){
    printf "Exiting:%s\n" "$1"
    exit -1
}

source customers.sh

peer lifecycle chaincode package ms-chaincode.tar.gz --lang golang --path ./chaincode --label ms_0 | peer lifecycle chaincode install ms-chaincode.tar.gz