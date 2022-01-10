/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	mspID         = "Org1MSP"
	appUser       = "alice"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "medicinecontract"
)

func main() {
	log.Println("============ Application starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("\nError setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("\nFailed to create wallet: %v", err)
	}

	if !wallet.Exists(appUser) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("\nFailed to populate wallet contents: %v", err)
		}
	} else {
		log.Println("============ Sucessfully populated wallet ============")
	}

	ccpPath := filepath.Join(
		"..",
		"configuration",
		"gateway",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, appUser),
	)
	if err != nil {
		log.Fatalf("\nFailed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("\nFailed to get network: %v", err)
	}

	contract := network.GetContract(chaincodeName)

	// Initiliase the ledger with mock data.
	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of medical supply on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	// // Handling request from a customer for a certain medicine.
	// log.Println("--> Submit Transaction: Request, function sends request for medicine.")
	// result, err = contract.SubmitTransaction("Request", "Aspirin", "00001", "Alice")
	// if err != nil {
	// 	log.Fatalf("\nFailed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

	// // Handling cancelling request from a customer for a certain medicine.
	// log.Println("--> Submit Transaction: CancelRequest, function sends request for medicine.")
	// result, err = contract.SubmitTransaction("CancelRequest", "Aspirin", "00001", "Alice")
	// if err != nil {
	// 	log.Fatalf("\nFailed to Submit transaction: %v", err)
	// }
	// log.Println(string(result))

	// // Handling user wanting to check his/her transaction history.
	// log.Println("--> Submit Transaction: CheckUserHistory, function shows history.")
	// result, err = contract.SubmitTransaction("CheckUserHistory", "Alice")
	// if err != nil {
	// 	log.Fatalf("\nFailed to Submit transaction: %v", err)
	// }
	// printArray(result)

	// // Handling user wanting to see all available medicine matching the medicine name.
	// log.Println("--> Submit Transaction: SearchMedicineByName, function shows available medicine matching the medicine name.")
	// result, err = contract.SubmitTransaction("SearchMedicineByName", "Zestril")
	// if err != nil {
	// 	log.Fatalf("\nFailed to Submit transaction: %v", err)
	// }
	// printArray(result)

	// Handling user wanting to see all available medicine.
	log.Println("--> Submit Transaction: CheckAvailableMedicine, function shows all available medicine.")
	result, err = contract.SubmitTransaction("CheckAvailableMedicine")
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	printArray(result)

	log.Println("\n============ Application ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return fmt.Errorf("Keystore folder should contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	return wallet.Put(appUser, identity)
}

func printArray(result []byte) {
	if len(result) > 0 {
		log.Println(string(result))
	} else {
		log.Println("No transactions found on ledger.")
	}
}
