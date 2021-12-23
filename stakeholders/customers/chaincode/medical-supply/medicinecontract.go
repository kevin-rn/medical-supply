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
