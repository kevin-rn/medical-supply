package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/go-tpm/tpm2"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var (
	tpmSrkHandle        = 0x800000ff
	tpmPcrSelection     = tpm2.PCRSelection{Hash: tpm2.AlgSHA256, PCRs: []int{}}
	tpmDefaultKeyParams = tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagDecrypt | tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin | tpm2.FlagNoDA,
		AuthPolicy: []byte{},
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:  tpm2.AlgNull,
				Mode: tpm2.AlgUnknown,
			},
			KeyBits: 2048,
		},
	}
)

func main() {
	fmt.Println("-Start TPM Code-")
	rwc, tpmError := tpm2.OpenTPM()
	if tpmError != nil {
		fmt.Println("TPM couldn't be opened. \n", tpmError)
		return
	}
	defer rwc.Close()

	// Prints out 8 random bytes from the TPM
	randByte, byteError := tpm2.GetRandom(rwc, 8)
	if byteError != nil {
		fmt.Println("TPM couldn't create 8 random bytes \n", byteError)
		return
	}
	fmt.Println("This prints out 8 random bytes from the TPM: \n", randByte)

	fmt.Println("-End TPM Code-")
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	medicines := []Medicine{
		Medicine{Name: "Aspirin", Disease: "Pain management", Expiration: "2022.05.09", Price: "$10", Holder: "MedStore"},
		Medicine{Name: "Vicodin", Disease: "Pain management", Expiration: "2022.07.01", Price: "$14", Holder: "MedStore"},
		Medicine{Name: "Synthroid", Disease: "Thyroid deficiency", Expiration: "2021.12.03", Price: "$11", Holder: "MedStore"},
		Medicine{Name: "Delasone", Disease: "Arthritis", Expiration: "2022.09.12", Price: "$5", Holder: "MedStore"},
		Medicine{Name: "Amoxil", Disease: "Bacterial infections", Expiration: "2022.07.08", Price: "$9", Holder: "MedStore"},
		Medicine{Name: "Neurontin", Disease: "Seizures", Expiration: "2022.03.25", Price: "$13", Holder: "MedStore"},
		Medicine{Name: "Zestril", Disease: "Blood pressure", Expiration: "2022.03.11", Price: "$7", Holder: "MedStore"},
		Medicine{Name: "Lipitor", Disease: "High cholesterol", Expiration: "2022.01.06", Price: "$12", Holder: "MedStore"},
		Medicine{Name: "Glucophage", Disease: "Type 2 diabetes", Expiration: "2022.04.24", Price: "$8", Holder: "MedStore"},
		Medicine{Name: "Zofran", Disease: "Nausea", Expiration: "2022.02.04", Price: "$13", Holder: "MedStore"},
		Medicine{Name: "Ibuprofen", Disease: "Fever", Expiration: "2022.02.28", Price: "$12", Holder: "MedStore"},
	}
	for i, med := range medicines {
		medAsBytes, _ := json.Marshal(med)
		err := ctx.GetStub().PutState("MED"+strconv.Itoa(i), medAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}
