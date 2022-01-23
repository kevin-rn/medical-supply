## For Monitoring using Hyperledger Explorer

![alt](../images/explorer.png?raw=true "Hyperledger Explorer")
1. Start test-network using networkDeploy.sh
2. Go to explorer folder: ```cd medical-supply/explorer```
3. Run: ```docker-compose up -d``` to start the Hyperledger Explorer 
4. Go to ```https://localhost:8080``` for the Hyperledger Explorer.   
For the login screen:    
username: exploreradmin   
password: exploreradminpw  
Note: These can be changed in the ```test-network.json``` file.
5. Run: ```docker-compose down -v``` to stop the Hyperledger Explorer
