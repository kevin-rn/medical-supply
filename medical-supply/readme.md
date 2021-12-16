## Start network:
```
cd fabric-samples/commercial-paper
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

```
cd medical-supply/stakeholders/customers
```

```
source customers.sh
```

```
peer lifecycle chaincode package cp.tar.gz --lang node --path ./chaincode --label cp_0
```

```
peer lifecycle chaincode install cp.tar.gz
```

```
export PACKAGE_ID= <id obtained from previous command>
```

```
peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name papercontract -v 0 --package-id $PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA
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


