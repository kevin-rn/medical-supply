package main

import (
	"fmt"
	"log"
	"time"

	medicalsupply "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/medical-supply"
)

// Temporary main method for measuring tpm latency
// Change first letter of functions to capital to make them global instead of private functions.
func main() {
	log.Println("==== Measuring latency tpm hash and key generation =====")
	for i := 1; i <= 10; i++ {
		start := time.Now()
		str, _ := medicalsupply.TpmHash("teststring")

		fmt.Printf("Hash: %s, time: %s \n", str, time.Since(start))
	}

	for i := 1; i <= 10; i++ {
		start := time.Now()
		str, _ := medicalsupply.TpmKey()

		fmt.Printf("Key generation: %s, time: %s \n", str, time.Since(start))
	}

}
