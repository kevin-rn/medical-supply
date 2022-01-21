package medicalsupply

import (
	"encoding/json"
	"strings"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/ledger-api"
)

// ListInterface - Functions which a medicinelist should have.
type ListInterface interface {
	AddMedicine(*MedicalSupply) error
	GetMedicine(string, string) (*MedicalSupply, error)
	GetAllMedicineByName(string) ([]*MedicalSupply, error)
	GetAllMedicine() ([]*MedicalSupply, error)
	UpdateMedicine(*MedicalSupply) error
	DeleteMedicine(string, string) error
	AddTPMAuth(*TPMAuth) error
	ExistsTPMAuth(string) bool
	VerifyTPMAuth(string, string) (bool, error)
}

type list struct {
	statelist ledgerapi.StateListInterface
}

// AddMedicine - Adding medicine to the statelist.
func (msl *list) AddMedicine(medicine *MedicalSupply) error {
	return msl.statelist.AddState(medicine)
}

// GetMedicine - Retrieves medicine from the statelist.
func (msl *list) GetMedicine(medName string, medNumber string) (*MedicalSupply, error) {
	ms := new(MedicalSupply)

	// Set to lower case
	medName = strings.ToLower(medName)

	// Use composite key to retrieve the medicine.
	err := msl.statelist.GetState(CreateMedicalKey(medName, medNumber), ms, "medicalsupply")
	if err != nil {
		return nil, err
	}
	return ms, nil
}

// GetAllMedicineByName - Retrieves all medicine matching the medicine name from the statelist.
func (msl *list) GetAllMedicineByName(medName string) ([]*MedicalSupply, error) {
	// Set to lower case
	medName = strings.ToLower(medName)

	// GetAllStatesByPartialKey returns an iterator
	data, err := msl.statelist.GetAllStatesByPartialKey(medName)
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

// GetAllMedicine - Retrieves all medicine from the statelist.
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

// UpdateMedicine - Update medicine (MedicalSupply object) on the statelist.
func (msl *list) UpdateMedicine(medicine *MedicalSupply) error {
	return msl.statelist.UpdateState(medicine)
}

// GetMedicine - Retrieves medicine from the statelist.
func (msl *list) DeleteMedicine(medName string, medNumber string) error {
	// Set to lower case
	medName = strings.ToLower(medName)
	return msl.statelist.DeleteState(CreateMedicalKey(medName, medNumber))
}

//-------------------------------------------------------//

// AddTPMAuth - Add tpm authentication to the ledger.
func (msl *list) AddTPMAuth(auth *TPMAuth) error {
	return msl.statelist.AddState(auth)
}

// GetTPMAuth - Check if TPM auth exists on the ledger.
func (msl *list) ExistsTPMAuth(holder string) bool {
	auth := new(TPMAuth)

	// Use composite key to retrieve the medicine.
	err := msl.statelist.GetState(createTPMledgerKey(holder), auth, "tpmauth")
	return err != nil
}

// VerifyTPMAuth - Check if TPM auth exists and verify the provided tpm key matches.
func (msl *list) VerifyTPMAuth(holder string, tpmkey string) (bool, error) {
	auth := new(TPMAuth)
	// Use composite key to retrieve the medicine.
	err := msl.statelist.GetState(createTPMledgerKey(holder), auth, "tpmauth")
	if err != nil {
		return false, err
	}
	// Verify if tpm key matches
	check := auth.TPMKey == tpmkey
	return check, nil
}

//-------------------------------------------------------//

// newList - Create new statelist.
func newList(ctx TransactionContextInterface) *list {
	statelist := new(ledgerapi.StateList)
	statelist.Ctx = ctx
	statelist.Name = "org.medstore.medicalsupplylist"
	statelist.DeserializeJSON = func(bytes []byte, state ledgerapi.StateInterface) error {
		return DeserializeJSON(bytes, state.(*MedicalSupply))
	}
	statelist.DeserializeTPM = func(bytes []byte, state ledgerapi.StateInterface) error {
		return DeserializeTPM(bytes, state.(*TPMAuth))
	}
	list := new(list)
	list.statelist = statelist
	return list
}
