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
	mspID         = "Org2MSP"
	appUser       = "bob"
	peerEndpoint  = "localhost:9051"
	gatewayPeer   = "peer0.org2.example.com"
	channelName   = "mychannel"
	chaincodeName = "medicinecontract"
)

func main() {
	log.Println("============ application starts ============")

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
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org2.example.com",
		"connection-org2.yaml",
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
	result, initerr := contract.SubmitTransaction("InitLedger")
	if initerr != nil {
		log.Fatalf("Failed to Submit transaction: %v", initerr)
	}
	log.Println(string(result))

	// Handling when regulators issue a new medicine (add to the ledger).
	log.Println("--> Submit Transaction: Issue, function sends issue for medicine.")
	result, issueerr := contract.SubmitTransaction("Issue", "Aspirin", "00012", "Pain management", "2022.05.09", "$10")
	if issueerr != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", issueerr)
	}
	log.Println(string(result))

	// Handling checking the entire transaction history.
	log.Println("--> Submit Transaction: CheckHistory, function shows history.")
	result, historyerr := contract.SubmitTransaction("CheckHistory")
	if historyerr != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", historyerr)
	}
	if len(result) > 0 {
		log.Println(string(result))
	} else {
		log.Println("Ledger has no transaction history.")
	}

	// Approves a medicine (changes its state from REQUESTED to SEND)
	log.Println("--> Submit Transaction: Approve, function that approves medicine.")
	result, approveerr := contract.SubmitTransaction("Approve", "Aspirin", "00001")
	if approveerr != nil {
		log.Fatalf("\nFailed to Submit transaction: %v", approveerr)
	}
	log.Println(string(result))

	log.Println("\n============ application ends ============")
}

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
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspID, string(cert), string(key))

	return wallet.Put(appUser, identity)
}
