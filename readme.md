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
cd medical-supply/customers
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
cd medical-supply/customers

source customers.sh
```

### For Regulators:
```
cd medical-supply/regulators
```

```
source regulators.sh
```

```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name medicinecontract -v 0 --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```
__________________________
## Application:
```

cd medical-supply/customers/application/


golang go: go run app.go
```
__________________________
## Clean (stop) Network:
```
cd fabric-samples/medical-supply
./network-clean.sh
```


