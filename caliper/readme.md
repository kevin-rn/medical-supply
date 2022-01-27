# Hyperledger Caliper
Custom benchmark tests have been implemented which can be found in the workload folder. The benchmarks/medicalSupplyBenchmark.yaml file defines which tests to run and the settings to use. The networks/networkConfig.yaml file defines which channel and user to use to connect to the fabric network (test network).

## For Running benchmark tests using Hyperledger Caliper
![alt](../images/caliper.png?raw=true "Hyperledger Caliper")
1. Install npm and run ```npm install``` inside caliper folder
2. Start test-network using ```source networkDeploy.sh```.
3. run ```source setup.sh``` for both customers and regulators to deploy the chaincode to test for.
4. run ```npx caliper bind --caliper-bind-sut fabric:2.2``` to bind hyperledger caliper to hyperledger fabric. Note: fabric version 2.3 did not work at the time of this project.
5. run ```npx caliper launch manager --caliper-fabric-gateway-enabled``` to start running the tests.
