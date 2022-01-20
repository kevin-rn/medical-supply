## Installing and setup:
1. install [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/latest/getting_started.html)  
2. Place this medical-supply repository in the fabric-samples repository  
3. Make sure to run ```go mod vendor``` in each folder with a go.mod file to get all the dependencies

_________________________
## Starting and stopping the network:
Go to ```cd medical-supply```  
Starting the network: ```source networkDeploy.sh```  
Stopping the network: ```source networkClean.sh```  

_________________________
## Deploying the chaincode (smart contract):

### For customers:
```
cd medical-supply/customers

source customers.sh
```

### For Regulators:
```
cd medical-supply/regulators

source regulators.sh
```

__________________________
## Application:
```
cd medical-supply/customers/application/
go run app.go
```

__________________________
## For Monitoring docker

To see the Fabric nodes running on the local machine: ```docker ps```
To view the network: ```docker network inspect fabric_test```

peer0.org1.example.com will be used for the Customers  
peer0.org2.example.com will be used for the Regulators


## For Monitoring as either customer or regulator
Go to their respective folder
``` cd medical-supply/customers ``` or ``` cd medical-supply/regulators ```

To show output from the Docker containers:
```./configuration/cli/monitordocker.sh fabric_test``` or alternatively if port number doesn't work: ```./monitordocker.sh fabric_test <port_number>```

## For Monitoring using Hyperledger Explorer

![alt](images/explorer.png?raw=true "Hyperledger Explorer")
1. Start test-network using networkDeploy.sh
2. Go to explorer folder: ```cd medical-supply/explorer```
3. Run: ```docker-compose up -d``` to start the Hyperledger Explorer 
4. Go to ```https://localhost:8080``` for the Hyperledger Explorer.   
For the login screen:    
username: exploreradmin   
password: exploreradminpw  
Note: These can be changed in the ```test-network.json``` file.
5. Run: ```docker-compose down -v``` to stop the Hyperledger Explorer

## For Running benchmark tests using Hyperledger Caliper
1. Start test-network using ```source networkDeploy.sh```.
2. run ```source customersSetup.sh``` and ```source regulatorsSetup.sh``` just like when deploying the chaincode.
3. run ```npx caliper bind --caliper-bind-sut fabric:2.2``` to bind hyperledger caliper to hyperledger fabric. Note: fabric version 2.3 did not work at the time of this project.
4. run ```npx caliper launch manager --caliper-fabric-gateway-enabled``` to start running the tests.