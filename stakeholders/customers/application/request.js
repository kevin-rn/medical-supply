'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');
const MedicalSupply = require('../chaincode-go/medical-supply/medicine.go');


// Main function
async function main() {
    // Connects with the wallet
    const wallet = await Wallets.newFileSystemWallet('../identity/user/alice/wallet');
    // Create new gateway
    const gateway = new Gateway();

    try {
        // Specify user who wants to request the medicine.
        const userName = 'alice';

        // Load connection profile and set up all options
        let connectionProfile = yaml.safeload(fs.readFileSync('../configuration/gateway/connection-org1.yaml', 'utf8'));
        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: {enabled: true, asLocalhost: true}
        };

        console.log('Connect to the Fabric gateway.');
        await gateway.connect(connectionProfile, connectionOptions);

        console.log('Use network channel: mychannel.');
        const network = await gateway.getNetwork('mychannel');

        console.log('Use org.medstore.medicalsupply smart contract');
        const contract = await network.getContract('medstore', 'org.medstore.medicalsupply')


        // Request the medicine:
        console.log('Submit the medical supply request transaction');
        // first string specifies type of transaction, so in this case a request.
        const requestResponse = await contract.submitTransaction('request', 'Aspirin', '00001', 'Pain management', '2022.05.09', '$10', 'MedStore')

        // process response
        console.log('Process request transaction response.'+requestResponse);
        console.log(`${medicine.holder} medical supply: ${medicine.medName} successfully requested`);
        let medicine = MedicalSupply.fromBuffer(requestResponse)

        console.log('Transaction complete.');

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {
            // Disconnect from the gateway
            console.log('Disconnect from Fabric gateway.');
            gateway.disconnect();
    }
}


// Calls main function
main().then(() => {

    console.log('Issue program complete.');

}).catch((e) => {

    console.log('Issue program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});