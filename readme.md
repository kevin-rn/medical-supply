## Error

kevin@vm:~/go/src/github.com/hyperledger/fabric-samples/medical-supply/stakeholders/customers/application$ node request.js
/home/kevin/go/src/github.com/hyperledger/fabric-samples/medical-supply/stakeholders/customers/chaincode-go/medical-supply/medicine.go:1
package medicalsupply
        ^^^^^^^^^^^^^

SyntaxError: Unexpected identifier
    at Object.compileFunction (node:vm:352:18)
    at wrapSafe (node:internal/modules/cjs/loader:1026:15)
    at Module._compile (node:internal/modules/cjs/loader:1061:27)
    at Object.Module._extensions..js (node:internal/modules/cjs/loader:1149:10)
    at Module.load (node:internal/modules/cjs/loader:975:32)
    at Function.Module._load (node:internal/modules/cjs/loader:822:12)
    at Module.require (node:internal/modules/cjs/loader:999:19)
    at require (node:internal/modules/cjs/helpers:102:18)
    at Object.<anonymous> (/home/kevin/go/src/github.com/hyperledger/fabric-samples/medical-supply/stakeholders/customers/application/request.js:7:23)
    at Module._compile (node:internal/modules/cjs/loader:1097:14)

Node.js v17.3.0

## Start network:
```
cd fabric-samples/medical-supply
./network-starter.sh
```
##### To see the Fabric nodes running on the local machine:
```
docker ps
```
##### To view the network:
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

```
source customers.sh
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


