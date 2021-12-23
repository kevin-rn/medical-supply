## Start network:
```
cd fabric-samples/medical-supply
./network-starter.sh
```
#### To see the Fabric nodes running on the local machine:
```
docker ps
```
#### To view the network:
```
docker network inspect fabric_test
```

peer0.org1.example.com will be used for the Customers
peer0.org2.example.com will be used for the Regulators

__________________________
## Monitor network as Customers or Regulators:
Go to the folder, for example:
```
cd medical-supply/stakeholders/customers
```

To show output from the Docker containers:
```
./configuration/cli/monitordocker.sh fabric_test
```
alternatively if port number doesn't work:
```
./monitordocker.sh fabric_test <port_number>
```
__________________________
## Deploying the chaincode (smart contract):

### For customers:
```
cd medical-supply/stakeholders/customers
```
#### Sets certain environment variables in command window (administrator) in order to use the correct set of peer binaries, send commands to the address of the organisation peer, and sign requests with the correct cryptographic material.
```
source customers.sh
```
#### Package the smart contract into a chaincode.
```
peer lifecycle chaincode package cp.tar.gz --lang node --path ./chaincode --label cp_0
```
#### Install chaincode
```
peer lifecycle chaincode install cp.tar.gz
```

#### Query the installed chaincode to get the package_id (same as when installing the chaincode)
```
peer lifecycle chaincode queryinstalled
```
#### Sets package_id as environmental variable.
```
export PACKAGE_ID= <id obtained from previous command>
```
##### Approve chaincode for the organisation
```
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name medicinecontract -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA
```
### For Regulators:
```
cd medical-supply/stakeholders/regulators
```

```
source regulators.sh
```

```
peer lifecycle chaincode package cp.tar.gz --lang node --path ./chaincode-go --label cp_0
```

```
peer lifecycle chaincode install cp.tar.gz
```

```
peer lifecycle chaincode queryinstalled
```

```
export PACKAGE_ID= <id obtained from previous command>
```

```
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name medicinecontract -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA
```

```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name medicinecontract -v 0 --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```
__________________________
## Application:
```
cd medical-supply/stakeholders/customers/application/

npm install
```

## Wallet
```
node enrollUserAlice.js
ls ../identity/user/alice/wallet/
cat ../identity/user/alice/wallet/*
```

## Run an application:
```
node <application name>

for example: node issue.js
```
__________________________
## Clean (stop) Network:
```
cd fabric-samples/medical-supply
./network-clean.sh
```


