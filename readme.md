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

source setup.sh
```

### For Regulators:
```
cd medical-supply/regulators

source setup.sh
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

## possible errors:
problem: Got permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock
run: sudo chmod 666 /var/run/docker.sock

problem: permission denied to /dev/tpm0 or /dev/tpmrm0
run: sudo chown <username> /dev/tpm0

replace - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock in 

lsof -n -i :9999
kill -9 <PID>