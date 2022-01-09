package medicalsupply

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("- Contract Instantiated -")
}

// InitLedger adds a base set of cars to the ledger
func (s *Contract) InitLedger(ctx TransactionContextInterface) error {
	medicines := []MedicalSupply{
		{MedName: "Aspirin", MedNumber: "00001", Disease: "Pain management", Expiration: "2022.05.09", Price: "$10", Holder: "MedStore"},
		{MedName: "Vicodin", MedNumber: "00002", Disease: "Pain management", Expiration: "2022.07.01", Price: "$14", Holder: "MedStore"},
		{MedName: "Synthroid", MedNumber: "00003", Disease: "Thyroid deficiency", Expiration: "2021.12.03", Price: "$11", Holder: "MedStore"},
		{MedName: "Delasone", MedNumber: "00004", Disease: "Arthritis", Expiration: "2022.09.12", Price: "$5", Holder: "MedStore"},
		{MedName: "Amoxil", MedNumber: "00005", Disease: "Bacterial infections", Expiration: "2022.07.08", Price: "$9", Holder: "MedStore"},
		{MedName: "Neurontin", MedNumber: "00006", Disease: "Seizures", Expiration: "2022.03.25", Price: "$13", Holder: "MedStore"},
		{MedName: "Zestril", MedNumber: "00007", Disease: "Blood pressure", Expiration: "2022.03.11", Price: "$7", Holder: "MedStore"},
		{MedName: "Lipitor", MedNumber: "00008", Disease: "High cholesterol", Expiration: "2022.01.06", Price: "$12", Holder: "MedStore"},
		{MedName: "Glucophage", MedNumber: "00009", Disease: "Type 2 diabetes", Expiration: "2022.04.24", Price: "$8", Holder: "MedStore"},
		{MedName: "Zofran", MedNumber: "00010", Disease: "Fever", Expiration: "2022.02.04", Price: "$13", Holder: "MedStore"},
		{MedName: "Ibuprofen", MedNumber: "00011", Disease: "Fever", Expiration: "2022.02.28", Price: "$12", Holder: "MedStore"},
	}
	for _, med := range medicines {

		err := ctx.GetMedicineList().UpdateMedicine(&med)

		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Function for handling issued medicine (Regulators)
func (c *Contract) Issue(ctx TransactionContextInterface, medname string, mednumber string,
	disease string, expiration string, price string) (*MedicalSupply, error) {

	// checkSumStr := fmt.Sprintf(medname, mednumber, disease, expiration, price, "MedStore")
	// checksum, _, tpmError := tpmHash(checkSumStr)

	// if tpmError != nil {
	// 	return nil, fmt.Errorf("Can't open TPM: %s", tpmError)
	// }
	checksum := "05010"

	// Create MedicalSupply object and set the State to Requested
	medicine := MedicalSupply{
		CheckSum:   checksum,
		MedName:    medname,
		MedNumber:  mednumber,
		Disease:    disease,
		Expiration: expiration,
		Price:      price,
		Holder:     "MedStore",
	}
	// Set medicine status to available
	medicine.SetAvailable()

	// Add the medicine to the medicine list to keep track
	err := ctx.GetMedicineList().AddMedicine(&medicine)
	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

// Function for handling requested medicine (Customers)
func (c *Contract) Request(ctx TransactionContextInterface, medname string, mednumber string, customer string) (*MedicalSupply, error) {
	medicine, err := ctx.GetMedicineList().GetMedicine(medname, mednumber)
	if err != nil {
		return nil, err
	}

	if medicine.Holder != "MedStore" {
		return nil, fmt.Errorf("medicine %s:%s has already been bought", medname, mednumber)
	}

	if medicine.IsAvailable() {
		medicine.SetRequested()
	} else {
		return nil, fmt.Errorf("medicine %s:%s is currently not available at MedStore", medname, mednumber)
	}

	if !medicine.IsRequested() {
		return nil, fmt.Errorf("medicine %s:%s is not requested. current state = %s", medname, mednumber, medicine.GetState())
	}

	medicine.Holder = customer

	err = ctx.GetMedicineList().UpdateMedicine(medicine)

	if err != nil {
		return nil, err
	}

	return medicine, nil

}

// Function for getting all Medicine
func (c *Contract) CheckHistory(ctx TransactionContextInterface, holder string) ([]*MedicalSupply, error) {
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, err
	}

	return medicinelist, nil
}

// Function for handling approving the medicine by marking it with Send (Regulators)
func (c *Contract) Approve(ctx TransactionContextInterface, medName string, medNumber string) (*MedicalSupply, error) {
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)

	if err != nil {
		return nil, err
	}

	if medicine.IsRequested() {
		medicine.SetSent()
	}

	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, err
	}

	return medicine, nil

}
