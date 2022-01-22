// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/hyperledger/fabric-chaincode-go/shim"
// 	"github.com/hyperledger/fabric-contract-api-go/contractapi"
// 	medicalsupply "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/medical-supply"
// )

// // Main method of the Chaincode (Smart contract).
// // Initialises all important information for when chaincode is packaged and installed on a channel.
// func main() {
// 	ccid := os.Getenv("CHAINCODE_ID")
// 	address := os.Getenv("CHAINCODE_SERVER_ADDRESS")

// 	contract := new(medicalsupply.Contract)
// 	contract.TransactionContextHandler = new(medicalsupply.TransactionContext)
// 	contract.Name = "org.medstore.medicalsupply"
// 	contract.Info.Version = "0.0.1"

// 	chaincode, err := contractapi.NewChaincode(contract)

// 	if err != nil {
// 		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
// 	}

// 	chaincode.Info.Title = "MedicalSupplyChaincode"
// 	chaincode.Info.Version = "0.0.1"

// 	server := &shim.ChaincodeServer{
// 		CCID:    ccid,
// 		Address: address,
// 		CC:      chaincode,
// 		TLSProps: shim.TLSProperties{
// 			Disabled: true,
// 		},
// 	}
// 	err = server.Start()
// 	if err != nil {
// 		fmt.Printf("Error starting medical supply external chaincode: %s", err)
// 	}

// }

package main

import (
	"fmt"

	"github.com/google/go-tpm/tpm2"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// For samples check https://github.com/hyperledger-archives/education/blob/master/LFS171x/fabric-material/chaincode/sample-chaincode.go

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (s *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("The function was invoked\n")
	// init code
	rwc, err := tpm2.OpenTPM("/dev/tpmrm0")
	if err != nil {
		shim.Error(fmt.Sprintf("could not open tpm:, %v", err))
	}
	res, err := tpm2.GetRandom(rwc, 2)
	if err != nil {
		shim.Error(fmt.Sprintf("could not get random numbers: %v", err))
	}
	fmt.Printf("Yoo it works res is %v\n", res)
	return shim.Success(res)
}

//NOTE - parameters such as ccid and endpoint information are hard coded here for illustration. This can be passed in in a variety of standard ways
func main() {
	// rwc, err_ := tpm2.OpenTPM("/dev/tpmrm0")
	// if err_ != nil {
	// 	fmt.Print("oeps eerste")
	// 	// shim.Error(fmt.Errorf("could not open tpm:, %v", err_).Error())
	// }
	// res, err_ := tpm2.GetRandom(rwc, 2)
	// if err_ != nil {
	// 	fmt.Print("oeps tweede")
	// 	// shim.Error(fmt.Errorf("could not get random numbers: %v", err_).Error())
	// }
	// fmt.Printf("%v", res)
	//The ccid is assigned to the chaincode on install (using the “peer lifecycle chaincode install <package>” command) for instance
	ccid := "mycc:724bf2be51e5a9e98d79d15d482d4fb1666af022c7f1368de18dec355d839da8"

	server := &shim.ChaincodeServer{
		CCID:    ccid,
		Address: "localhost:9999",
		CC:      new(SimpleChaincode),
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}
	fmt.Printf("Starting chaincode %s at %s\n", server.CCID, server.Address)
	err := server.Start()
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
