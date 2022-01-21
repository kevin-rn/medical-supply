package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	medicalsupply "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/medical-supply"
)

// Main method of the Chaincode (Smart contract).
// Initialises all important information for when chaincode is packaged and installed on a channel.
func main() {
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

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
