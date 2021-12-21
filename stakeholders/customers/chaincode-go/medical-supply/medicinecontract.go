package medicalsupply

import (
	"flag"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	"github.com/google/go-tpm/tpm2"
)

var (
	srkTemplate = tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth | tpm2.FlagRestricted | tpm2.FlagDecrypt | tpm2.FlagNoDA,
		AuthPolicy: nil,
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
			KeyBits:    2048,
			ModulusRaw: make([]byte, 256),
		},
	}
	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket).")
	pcr     = flag.Int("pcr", -1, "PCR to seal data to. Must be within [0, 23].")
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("- Contract Instantiated -")
}

// Function for handling requested medicine
func (c *Contract) Request(ctx TransactionContextInterface, medname string, mednumber string,
	disease string, expiration string, price string, holder string) (*MedicalSupply, error) {

	// sealedHolder = seal(holder)
	// Create MedicalSupply object and set the State to Requested
	medicine := MedicalSupply{MedName: medname, MedNumber: mednumber, Disease: disease, Expiration: expiration, Price: price, Holder: holder}
	medicine.SetRequested()

	err := ctx.GetMedicineList().AddMedicine(&medicine)
	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

// Function for handling send medicine
func (c *Contract) Send(ctx TransactionContextInterface, medName string, medNumber string, oldHolder string, newHolder string) (*MedicalSupply, error) {
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)

	if err != nil {
		return nil, err
	}

	if medicine.Holder != oldHolder {
		return nil, fmt.Errorf("Medicine %s:%s is not owned by %s", medName, medNumber, oldHolder)
	}

	if medicine.IsRequested() {
		medicine.SetSent()
	}

	medicine.Holder = newHolder

	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, err
	}

	return medicine, nil

}

// func seal(string info) (result []byte, retError error) {
// 	// Connect to the TPM
// 	rwc, err := tpm2.OpenTPM(*tpmPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't open TPM %q: %v", tpmPath, err)
// 	}
// 	defer func() {
// 		if err := rwc.Close(); err != nil {
// 			retError = fmt.Errorf("%v\ncan't close TPM %q: %v", retError, tpmPath, err)
// 		}
// 	}()

// 	// Create the parent key against which to seal the data
// 	srkPassword := ""
// 	srkHandle, _, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", srkPassword, srkTemplate)
// 	if err != nil {
// 		return fmt.Errorf("can't create primary key: %v", err)
// 	}
// 	defer func() {
// 		if err := tpm2.FlushContext(rwc, srkHandle); err != nil {
// 			retErr = fmt.Errorf("%v\nunable to flush SRK handle %q: %v", retErr, srkHandle, err)
// 		}
// 	}()
// 	fmt.Printf("Created parent key with handle: 0x%x\n", srkHandle)

// 	// Note the value of the pcr against which we will seal the data
// 	pcrVal, err := tpm2.ReadPCR(rwc, pcr, tpm2.AlgSHA256)
// 	if err != nil {
// 		return fmt.Errorf("unable to read PCR: %v", err)
// 	}
// 	fmt.Printf("PCR %v value: 0x%x\n", pcr, pcrVal)

// 	// Get the authorization policy that will protect the data to be sealed
// 	objectPassword := "objectPassword"
// 	sessHandle, policy, err := policyPCRPasswordSession(rwc, pcr, objectPassword)
// 	if err != nil {
// 		return fmt.Errorf("unable to get policy: %v", err)
// 	}
// 	if err := tpm2.FlushContext(rwc, sessHandle); err != nil {
// 		return fmt.Errorf("unable to flush session: %v", err)
// 	}
// 	fmt.Printf("Created authorization policy: 0x%x\n", policy)

// 	// Seal the data to the parent key and the policy
// 	dataToSeal := []byte("secret")
// 	fmt.Printf("Data to be sealed: 0x%x\n", dataToSeal)
// 	privateArea, publicArea, err := tpm2.Seal(rwc, srkHandle, srkPassword, objectPassword, policy, dataToSeal)
// 	if err != nil {
// 		return fmt.Errorf("unable to seal data: %v", err)
// 	}
// 	fmt.Printf("Sealed data: 0x%x\n", privateArea)
// 	return publicArea
// }
