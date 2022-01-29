# Medical Supply - Let smart contract meet trusted hardware based attestation 
This Medical-Supply repository contains a prototype, implemented on [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/latest/whatis.html), of a Medical Supply Chain with TPM functions integrated. Here the goal is to show how the cryptographic functionalities of TPM 2.0 can be combined with the smart contract used in a blockchain application to provide more security. For utilising TPM 2.0, this repository relies on the [Go-TPM](https://github.com/google/go-tpm) repository. 

Medical-Supply is part of a bachelor thesis project at the Delft university of Technology. It should be noted that combining TPM 2.0 with the prototype was not succesful due to Docker not being able to access the '/dev/tpmrm0' folder as it is not part of its container environment. An attempt was made at using [Chaincode as external service](https://hyperledger-fabric.readthedocs.io/en/latest/cc_service.html) as solution due to it being able to access the local system instead of being installed on the peers in Docker. However this was not accomplished within the time frame of my project. The external chaincode setup can be found on the external-chaincode branch. The main branch will have the TPM functions hardcoded to return a string.

## Context of the Medical Supply Chain prototype
The innerworkings are similar to the [Commercial Paper](https://github.com/hyperledger/fabric-samples/tree/main/commercial-paper) example. A pharmacy network called MedStore is used (this is the test-network provided by Hyperledger Fabric). On this network, Medical Supply form the assets. There are two types of 'organisations' called Customers which can for example request medical supply and Regulators which can for example issue medical supply. Regulators form a more administrative role.

_________________________
## Installing and setup:
1. Install [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/latest/getting_started.html). Make sure to follow the Prerequisites steps.
2. Clone this 'medical-supply' repository and place it in the fabric-samples repository  
3. Run ```go mod vendor``` in each folder with a go.mod file (application and chaincode) to get all the dependencies

_________________________
## Running the prototype:
Go to ```medical-supply```  folder in the terminal.  
Starting the network: 
```
medical-supply$ source networkDeploy.sh
```  
Install the chaincode (smart contract) on customers first:
```
medical-supply/customers$ source setup.sh
```

Install the chaincode (smart contract) on regulators:
```
medical-supply/regulators$ source setup.sh
```
For running the application, Go to either ```customers/application``` or ```regulators/application```. Now run:
```
../application$ go run app.go
```

Stopping the network: 
```
medical-supply$ source networkClean.sh
```  
__________________________
## For Monitoring Docker through the terminal:

To see the Fabric nodes running on the local machine: ```docker ps```  
To view the network: ```docker network inspect fabric_test```  

peer0.org1.example.com will be used for the Customers  
peer0.org2.example.com will be used for the Regulators

__________________________
## For Monitoring as either Customers or Regulators through the terminal:
Go to their respective folder: 
``` medical-supply/customers ``` or ``` medical-supply/regulators ```

To show output from the Docker containers:
```
./configuration/cli/monitordocker.sh fabric_test
``` 
or alternatively if port number doesn't work: 
```
./monitordocker.sh fabric_test <port_number>
```

__________________________
## Additional integrated tools:
- [Hyperledger Caliper](https://github.com/hyperledger/caliper/) for doing performance analysis. For running it see the readme in the caliper folder.  
 - [Hyperledger Explorer](https://github.com/hyperledger/blockchain-explorer) for a userfriendly web application. For running it see the readme in the explorer folder. 
__________________________
## Possible solutions for problems on Linux:
problem: Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock  
run: sudo chmod 666 /var/run/docker.sock

problem: permission denied to /dev/tpm0 or /dev/tpmrm0  
run: sudo chown <username> /dev/tpm0
