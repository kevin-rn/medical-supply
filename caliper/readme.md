## For Running benchmark tests using Hyperledger Caliper
![alt](../images/caliper.png?raw=true "Hyperledger Caliper")
1. Install npm and run ```npm install``` inside caliper folder
2. Start test-network using ```source networkDeploy.sh```.
3. run ```source setup.sh``` and ```source setup.sh``` just like when deploying the chaincode.
4. run ```npx caliper bind --caliper-bind-sut fabric:2.2``` to bind hyperledger caliper to hyperledger fabric. Note: fabric version 2.3 did not work at the time of this project.
5. run ```npx caliper launch manager --caliper-fabric-gateway-enabled``` to start running the tests.
