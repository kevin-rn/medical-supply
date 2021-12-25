package medicalsupply

import (
	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/stakeholders/customers/chaincode/ledger-api"
)

type ListInterface interface {
	AddMedicine(*MedicalSupply) error
	GetMedicine(string, string) (*MedicalSupply, error)
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
