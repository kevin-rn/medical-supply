package main

import (
	"fmt"
	"os"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	medicalsupply "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/medical-supply"
)

// Main method of the Chaincode (Smart contract).
// Initialises all important information for when chaincode is packaged and installed on a channel.
func main() {
	ccid := os.Getenv("CHAINCODE_ID")
	address := os.Getenv("CHAINCODE_SERVER_ADDRESS")

	contract := new(medicalsupply.Contract)
	contract.TransactionContextHandler = new(medicalsupply.TransactionContext)
	contract.Name = "org.medstore.medicalsupply"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "MedicalSupplyChaincode"
	chaincode.Info.Version = "0.0.1"

	server := &shim.ChaincodeServer{
		CCID:    ccid,
		Address: address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}
	err = server.Start()
	if err != nil {
		fmt.Printf("Error starting medical supply external chaincode: %s", err)
	}

}
