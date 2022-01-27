## Running Chaincode as external service
Instead of installing chaincode on a peer, the chaincode itself can run as its own service where it can access the local file system.

Instructions to run chaincode as external service (based on [Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io/en/latest/cc_service.html)).   
**Note: This has not succesfully run and therefore these instructions might not work**
1. Replace the `docker-compose-test-net.yaml` in fabric-samples/test-network/docker with the one in medical-supply/external-chaincode.  
**Note: The settings inside this file for running chaincode as external service do not work**.
2. Adjust the networkClean.sh and networkDeploy.sh to point to the local `medical-supply/config` folder instead of the `fabric-samples/config` folder.
3. Run the network using `source networkDeploy.sh`
4. Go to the external-chaincode folder in the terminal.
5. Run `tar cfz code.tar.gz connection.json`.
6. Run `tar cfz ms-external-chaincode.tgz metadata.json code.tar.gz`.
7. ~~Run `source setup-external.sh` in customers folder and afterwards the regulators folder~~.   
Go to first Customers, export global environment variables for organisation 1 and afterwards run   
`peer lifecycle chaincode install ../external-chaincode/ms-external-chaincode.tgz`.  
Do the same for Regulators only export global enviroment variables for organisation 2.
8. Go back to the external-chaincode folder and run `go run main.go` to start external chaincode as a service
9. Stopping the network using `source networkClean.sh`


