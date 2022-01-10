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

// InitLedger adds a base set of medicine (MedicalSupply) to the ledger. [Regulators]
func (s *Contract) InitLedger(ctx TransactionContextInterface) error {
	// Create array of MedicalSupply objects.
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

	// For each medicine, set it's state to Available, calculate the checksum and update the ledger
	for _, med := range medicines {
		med.SetAvailable()
		// TODO: checksum
		err := ctx.GetMedicineList().UpdateMedicine(&med)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Function for handling issued medicine [Regulators]
func (c *Contract) Issue(ctx TransactionContextInterface, medname string, mednumber string,
	disease string, expiration string, price string) (*MedicalSupply, error) {

	// Calculate the checksum by using the hashfunction of the TPM.
	// checkSumStr := fmt.Sprintf(medname, mednumber, disease, expiration, price, "MedStore")
	// checksum, _, tpmError := tpmHash(checkSumStr)

	// if tpmError != nil {
	// 	return nil, fmt.Errorf("Can't open TPM: %s", tpmError)
	// }
	checksum := "05010"

	// Create MedicalSupply object.
	medicine := MedicalSupply{
		CheckSum:   checksum,
		MedName:    medname,
		MedNumber:  mednumber,
		Disease:    disease,
		Expiration: expiration,
		Price:      price,
		Holder:     "MedStore",
	}
	// Set state to AVAILABLE.
	medicine.SetAvailable()

	// Add the medicine to the ledger.
	err := ctx.GetMedicineList().AddMedicine(&medicine)
	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

// Function for handling requested medicine. [Customers]
func (c *Contract) Request(ctx TransactionContextInterface, medname string, mednumber string, customer string) (*MedicalSupply, error) {
	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medname, mednumber)
	if err != nil {
		return nil, err
	}

	// Verify that the current holder is MedStore, if that is not the case than the medicine has already been transferred to a different holder.
	if medicine.Holder != "MedStore" {
		return nil, fmt.Errorf("Medicine %s:%s has already been bought.", medname, mednumber)
	}

	// Verify that the current state is AVAILABLE, if so set to REQUESTED.
	if medicine.IsAvailable() {
		medicine.SetRequested()
	} else {
		return nil, fmt.Errorf("Medicine %s:%s is currently not available at MedStore.", medname, mednumber)
	}

	// Verify that change to REQUESTED state has succeeded.
	if !medicine.IsRequested() {
		return nil, fmt.Errorf("Medicine %s:%s is not requested. current state = %s.", medname, mednumber, medicine.GetState())
	}

	// Update medicine holder to be the customer instead of MedStore.
	medicine.Holder = customer
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, err
	}

	return medicine, nil
}

// Function for getting all Medicine. [Regulators]
func (c *Contract) CheckHistory(ctx TransactionContextInterface) ([]*MedicalSupply, error) {
	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, err
	}
	return medicinelist, nil
}

// Function for getting all available medicine. [Customers]
func (c *Contract) CheckAvailableMedicine(ctx TransactionContextInterface) ([]*MedicalSupply, error) {
	// Get all medicine from the ledger (There is currently no efficienter way to retrieve assets from the Ledger for certain fields).
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, err
	}

	// Loop through the list and check for AVAILABLE state.
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.IsAvailable() {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// Function for getting all requested medicine. [Regulators]
func (c *Contract) CheckRequestedMedicine(ctx TransactionContextInterface) ([]*MedicalSupply, error) {
	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, err
	}

	// Loop through the list and check for REQUESTED state.
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.IsRequested() {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// Function for getting all Medicine an User has ordered. [Customers]
func (c *Contract) CheckUserHistory(ctx TransactionContextInterface, holder string) ([]*MedicalSupply, error) {
	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, err
	}

	// Loop through the list and check for the user (holder).
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.Holder == holder {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// Function for handling approving the medicine by changing its state to SEND. [Regulators]
func (c *Contract) Approve(ctx TransactionContextInterface, medName string, medNumber string) (*MedicalSupply, error) {
	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, err
	}

	// Check if medicine state is REQUESTED, if so set it to SEND.
	if medicine.IsRequested() {
		medicine.SetSend()
	} else {
		return nil, fmt.Errorf("Cannot approve medicine that has not been requested.")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, err
	}

	return medicine, nil
}

// Function for handling disapproving the medicine by changing its state back to AVAILABLE. [Regulators]
func (c *Contract) Disapprove(ctx TransactionContextInterface, medName string, medNumber string) (*MedicalSupply, error) {
	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, err
	}

	// Check if medicine state is REQUESTED, if so set it to AVAILABLE and reset to holder to be MedStore.
	if medicine.IsRequested() {
		medicine.SetAvailable()
		medicine.Holder = "MedStore"
	} else {
		return nil, fmt.Errorf("Cannot disapprove medicine that has not been requested.")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, err
	}

	return medicine, nil
}
