package main

import (
	"fmt"

	"github.com/google/go-tpm/tpm2"
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
