/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	wallet := enrollUser()
	contract := connectToNetwork(wallet)

	log.Println("Choose number to invoke function: \n" +
		"1 - Request a medicine \n" +
		"2 - Cancel request \n" +
		"3 - Check User History \n" +
		"4 - Search Medicine by name \n" +
		"5 - Check available medicine")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch strings.ToLower(input) {
	case "1":
		request(contract, scanner)
	case "2":
		cancelrequest(contract, scanner)
	case "3":
		checkUserHistory(contract)
	case "4":
		searchMedicineByName(contract, scanner)
	case "5":
		checkAvailableMedicine(contract)
	default:
		log.Fatalf("\n Error: Function to invoke not found.")
	}

}

// Enrolls user as peer to the network
func enrollUser() *gateway.Wallet {
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

	return wallet
}

// Connects to the network channel and gets the smart contract to invoke functions on.
func connectToNetwork(wallet *gateway.Wallet) *gateway.Contract {
	ccpPath := filepath.Join("..", "configuration", "gateway", "connection-org1.yaml")
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
	return contract
}

// Create wallet and keystore folder for user to use.
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

// Helper function for pretty printing results to the terminal.
func prettyPrint(body []byte) {
	var indentedFormat bytes.Buffer
	error := json.Indent(&indentedFormat, body, "", "\t")
	if error != nil {
		log.Println("Error encountered Json parse error: ", error)
		// Print bytes normally without the pretty printed format
		log.Println(string(body))
		return
	}
	log.Println(string(indentedFormat.Bytes()))
}

// Helper function for printing array results.
func printArray(result []byte) {
	if len(result) > 0 {
		prettyPrint(result)
	} else {
		log.Println("No transactions found on ledger.")
	}
}

// Invokes function that puts a request for a certain medicine.
func request(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()

	log.Println("--> Submit Transaction: Request, function sends request for medicine.")
	result, err := contract.SubmitTransaction("Request", strings.ToLower(medName), medNumber, appUser)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Invokes function that cancels request for a certain medicine.
func cancelrequest(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()

	log.Println("--> Submit Transaction: CancelRequest, function sends request for medicine.")
	result, err := contract.SubmitTransaction("CancelRequest", medName, medNumber, appUser)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Invokes function that returns an user's transaction history.
func checkUserHistory(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: CheckUserHistory, function shows history.")
	result, err := contract.SubmitTransaction("CheckUserHistory", appUser)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Invokes function that returns all available medicine matching the medicine name.
func searchMedicineByName(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()

	log.Println("--> Submit Transaction: SearchMedicineByName, function shows available medicine matching the medicine name.")
	result, err := contract.SubmitTransaction("SearchMedicineByName", medName)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	printArray(result)
}

// Invokes function that returns all available medicine.
func checkAvailableMedicine(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: CheckAvailableMedicine, function shows all available medicine.")
	result, err := contract.SubmitTransaction("CheckAvailableMedicine")
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	printArray(result)
}
