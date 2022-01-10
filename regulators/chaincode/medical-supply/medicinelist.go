package medicalsupply

import (
	"encoding/json"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/regulators/chaincode/ledger-api"
)

// ListInterface functions which a medicinelist should have.
type ListInterface interface {
	AddMedicine(*MedicalSupply) error
	GetMedicine(string, string) (*MedicalSupply, error)
	GetAllMedicine() ([]*MedicalSupply, error)
	UpdateMedicine(*MedicalSupply) error
}

type list struct {
	statelist ledgerapi.StateListInterface
}

// Adding medicine to the statelist.
func (msl *list) AddMedicine(medicine *MedicalSupply) error {
	return msl.statelist.AddState(medicine)
}

// Retrieving medicine from the statelist.
func (msl *list) GetMedicine(medName string, medNumber string) (*MedicalSupply, error) {
	ms := new(MedicalSupply)

	// Use composite key to retrieve the medicine.
	err := msl.statelist.GetState(CreateMedicalKey(medName, medNumber), ms)
	if err != nil {
		return nil, err
	}
	return ms, nil
}

// Retrieving all medicine from the statelist.
func (msl *list) GetAllMedicine() ([]*MedicalSupply, error) {
	// GetAllStates returns an iterator
	data, err := msl.statelist.GetAllStates()
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// Use iterator to loop and return an array of all MedicalSupply objects.
	var medicines []*MedicalSupply
	for data.HasNext() {
		queryResponse, err := data.Next()
		if err != nil {
			return nil, err
		}

		var med MedicalSupply
		err = json.Unmarshal(queryResponse.Value, &med)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, &med)
	}

	return medicines, nil
}

// Update medicine (MedicalSupply object) on the statelist.
func (msl *list) UpdateMedicine(medicine *MedicalSupply) error {
	return msl.statelist.UpdateState(medicine)
}

// Create new statelist.
func newList(ctx TransactionContextInterface) *list {
	statelist := new(ledgerapi.StateList)
	statelist.Ctx = ctx
	statelist.Name = "org.medstore.medicalsupplylist"
	statelist.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*MedicalSupply))
	}

	list := new(list)
	list.statelist = statelist
	return list
}
