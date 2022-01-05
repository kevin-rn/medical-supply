package medicalsupply

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("- Contract Instantiated -")
}

// InitLedger adds a base set of cars to the ledger
func (s *Contract) InitLedger(ctx contractapi.TransactionContextInterface) error {
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
		{MedName: "Zofran", MedNumber: "00010", Disease: "Nausea", Expiration: "2022.02.04", Price: "$13", Holder: "MedStore"},
		{MedName: "Ibuprofen", MedNumber: "00011", Disease: "Fever", Expiration: "2022.02.28", Price: "$12", Holder: "MedStore"},
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

// Function for handling requested medicine
func (c *Contract) Request(ctx TransactionContextInterface, medname string, mednumber string,
	disease string, expiration string, price string, holder string) (*MedicalSupply, error) {

	// Create MedicalSupply object and set the State to Requested
	medicine := MedicalSupply{MedName: medname, MedNumber: mednumber, Disease: disease, Expiration: expiration, Price: price, Holder: holder}
	medicine.SetRequested()

	err := ctx.GetMedicineList().AddMedicine(&medicine)
	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

func (c *Contract) RequestRecommended(ctx TransactionContextInterface, disease string, expiration string, price string, holder string) (*MedicalSupply, error) {
	return nil, nil
}

func (c *Contract) CheckHistory(ctx TransactionContextInterface, holder string) (*MedicalSupply, error) {
	return nil, nil
}

func (c *Contract) CheckPending(ctx TransactionContextInterface, holder string) (*MedicalSupply, error) {
	return nil, nil
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
