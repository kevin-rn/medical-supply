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
	mspID         = "Org2MSP"
	appUser       = "bob"
	peerEndpoint  = "localhost:9051"
	gatewayPeer   = "peer0.org2.example.com"
	channelName   = "mychannel"
	chaincodeName = "medicinecontract"
)

func main() {
	wallet := enrollUser()
	contract := connectToNetwork(wallet)

	tpmkey, err := tpmKeyHandler(contract, "tpmkey.txt")
	if err != nil {
		log.Fatalf("Failed to generate TPM key: %v", err)
	}
	log.Printf("TPM Key used is: %v", tpmkey)

	log.Println("Choose number to invoke function: \n" +
		"1 - Initialise the ledger \n" +
		"2 - Check the entire ledger history \n" +
		"3 - Issue new medicine \n" +
		"4 - Change status of a medicine \n" +
		"5 - Change holder of medicine \n" +
		"6 - Check all requested medicine \n" +
		"7 - Approve request for medicine \n" +
		"8 - Reject request for medicine \n" +
		"9 - Delete medicine")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch strings.ToLower(input) {
	case "1":
		initLedger(contract)
	case "2":
		checkHistory(contract)
	case "3":
		issue(contract, scanner)
	case "4":
		changeStatusMedicine(contract, scanner)
	case "5":
		changeHolder(contract, scanner)
	case "6":
		checkRequestedMedicine(contract)
	case "7":
		approveRequest(contract, scanner)
	case "8":
		rejectRequest(contract, scanner)
	case "9":
		delete(contract, scanner)
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
	ccpPath := filepath.Join("..", "configuration", "gateway", "connection-org2.yaml")

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
		"org2.example.com",
		"users",
		"User1@org2.example.com",
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
		return fmt.Errorf("keystore folder should contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	return wallet.Put(appUser, identity)
}

// Reads tpm key from file, if no success then request for new key and store that.
func tpmKeyHandler(contract *gateway.Contract, filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		// Request tpm key from smart contract
		log.Println("--> Submit Transaction: TPMKeyGen, function requests for tpm generated key.")
		result, err := contract.SubmitTransaction("TPMKeyGen", appUser)
		if err != nil {
			log.Fatalf("\nFailed to Submit transaction: %v", err)
		}
		tpmkey := string(result)

		// Store tpmkey to file
		file, err := os.Create(filepath)
		if err != nil {
			return tpmkey, err
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		fmt.Fprintln(w, tpmkey)
		return tpmkey, w.Flush()

	} else {
		// Read key from file
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		line := scanner.Text()
		return line, scanner.Err()
	}
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

// Helper function for printing array results
func printArray(result []byte) {
	if len(result) > 0 {
		prettyPrint(result)
	} else {
		log.Println("No transactions found on ledger.")
	}
}

// Initiliase the ledger with mock data.
func initLedger(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of medical supply on the ledger")
	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
}

// Handling checking the entire transaction history.
func checkHistory(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: CheckHistory, function shows history.")
	result, err := contract.SubmitTransaction("CheckHistory")
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	printArray(result)
}

// Handling when regulators issue a new medicine (add to the ledger).
func issue(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00012):")
	scanner.Scan()
	medNumber := scanner.Text()
	log.Println("Disease (e.g. Pain management):")
	scanner.Scan()
	disease := scanner.Text()
	log.Println("Expiration date (e.g. 2022.05.09):")
	scanner.Scan()
	expirationDate := scanner.Text()
	log.Println("Price (e.g. $10):")
	scanner.Scan()
	price := scanner.Text()

	log.Println("--> Submit Transaction: Issue, function sends issue for medicine.")
	result, err := contract.SubmitTransaction("Issue", medName, medNumber, disease, expirationDate, price)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Changing status of medicine manually.
func changeStatusMedicine(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()
	log.Println("Medicine status (e.g. Available, Requested or Send):")
	scanner.Scan()
	status := scanner.Text()

	log.Println("--> Submit Transaction: ChangeStatusMedicine, function sends request for medicine.")
	result, err := contract.SubmitTransaction("ChangeStatusMedicine", medName, medNumber, status)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Changing holder of medicine manually.
func changeHolder(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()
	log.Println("Holder name (e.g. John):")
	scanner.Scan()
	holder := scanner.Text()

	log.Println("--> Submit Transaction: ChangeHolder, function sends request for medicine.")
	result, err := contract.SubmitTransaction("ChangeHolder", medName, medNumber, holder)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Handling regulators wanting to see all requested medicine matching the medicine name.
func checkRequestedMedicine(contract *gateway.Contract) {
	log.Println("--> Submit Transaction: CheckRequestedMedicine, function shows all requested medicine.")
	result, err := contract.SubmitTransaction("CheckRequestedMedicine")
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	printArray(result)
}

// Approves a medicine (changes its state from REQUESTED to SEND)
func approveRequest(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()

	log.Println("--> Submit Transaction: ApproveRequest, function that approves medicine.")
	result, err := contract.SubmitTransaction("ApproveRequest", medName, medNumber)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Rejecting a medicine (changes its state from REQUESTED to AVAILABLE)
func rejectRequest(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()

	log.Println("--> Submit Transaction: RejectRequest, function that approves medicine.")
	result, err := contract.SubmitTransaction("RejectRequest", medName, medNumber)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}

// Deletes a medicine
func delete(contract *gateway.Contract, scanner *bufio.Scanner) {
	log.Println("Medicine name (e.g. Aspirin):")
	scanner.Scan()
	medName := scanner.Text()
	log.Println("Medicine number (e.g. 00001):")
	scanner.Scan()
	medNumber := scanner.Text()

	log.Println("--> Submit Transaction: Delete, function that approves medicine.")
	result, err := contract.SubmitTransaction("Delete", medName, medNumber)
	if err != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", err)
	}
	prettyPrint(result)
}
