/*
 * SPDX-License-Identifier: Apache-2.0
 */

package ledgerapi

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	medicalsupply "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/medical-supply"
)

// MedicalSupply - Defines a medicine.
type MedicalSupply = medicalsupply.MedicalSupply

// StateListInterface functions that a state list
// should have
type StateListInterface interface {
	AddState(StateInterface) error
	GetState(string, StateInterface) error
	GetAllStates(string, string) ([]*MedicalSupply, error)
	UpdateState(StateInterface) error
}

// StateList useful for managing putting data in and out
// of the ledger. Implementation of StateListInterface
type StateList struct {
	Ctx         contractapi.TransactionContextInterface
	Name        string
	Deserialize func([]byte, StateInterface) error
}

// AddState puts state into world state
func (sl *StateList) AddState(state StateInterface) error {
	key, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, state.GetSplitKey())
	data, err := state.Serialize()

	if err != nil {
		return err
	}

	return sl.Ctx.GetStub().PutState(key, data)
}

// GetState returns state from world state. Unmarshalls the JSON
// into passed state. Key is the split key value used in Add/Update
// joined using a colon
func (sl *StateList) GetState(key string, state StateInterface) error {
	ledgerKey, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, SplitKey(key))
	data, err := sl.Ctx.GetStub().GetState(ledgerKey)

	if err != nil {
		return err
	} else if data == nil {
		return fmt.Errorf("no state found for %s", key)
	}

	return sl.Deserialize(data, state)
}

func (sl *StateList) GetAllStates(startkey string, endkey string) ([]*MedicalSupply, error) {
	ledgerStart, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, SplitKey(startkey))
	ledgerEnd, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, SplitKey(endkey))
	data, err := sl.Ctx.GetStub().GetStateByRange(ledgerStart, ledgerEnd)

	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, fmt.Errorf("no history on medicines found from MedStore")
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

// UpdateState puts state into world state. Same as AddState but
// separate as semantically different
func (sl *StateList) UpdateState(state StateInterface) error {
	return sl.AddState(state)
}
