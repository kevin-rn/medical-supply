package medicalsupply

import (
	"encoding/json"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/ledger-api"
)

type ListInterface interface {
	AddMedicine(*MedicalSupply) error
	GetMedicine(string, string) (*MedicalSupply, error)
	GetAllMedicine() ([]*MedicalSupply, error)
	UpdateMedicine(*MedicalSupply) error
}

type list struct {
	statelist ledgerapi.StateListInterface
}

func (msl *list) AddMedicine(medicine *MedicalSupply) error {
	return msl.statelist.AddState(medicine)
}

func (msl *list) GetMedicine(medName string, medNumber string) (*MedicalSupply, error) {
	ms := new(MedicalSupply)

	err := msl.statelist.GetState(CreateMedicalKey(medName, medNumber), ms)
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (msl *list) GetAllMedicine() ([]*MedicalSupply, error) {
	data, err := msl.statelist.GetAllStates()
	if err != nil {
		return nil, err
	}

	defer data.Close()

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

func (msl *list) UpdateMedicine(medicine *MedicalSupply) error {
	return msl.statelist.UpdateState(medicine)
}

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
