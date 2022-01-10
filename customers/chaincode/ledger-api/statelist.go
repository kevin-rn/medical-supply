/*
 * SPDX-License-Identifier: Apache-2.0
 */

package ledgerapi

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// StateListInterface functions that a state list should have.
type StateListInterface interface {
	AddState(StateInterface) error
	GetState(string, StateInterface) error
	GetAllStates() (shim.StateQueryIteratorInterface, error)
	UpdateState(StateInterface) error
}

// StateList useful for managing putting data in and out of the ledger.
// Implementation of StateListInterface.
type StateList struct {
	Ctx         contractapi.TransactionContextInterface
	Name        string
	Deserialize func([]byte, StateInterface) error
}

// AddState puts state into world state.
func (sl *StateList) AddState(state StateInterface) error {
	key, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, state.GetSplitKey())
	data, err := state.Serialize()

	if err != nil {
		return err
	}

	return sl.Ctx.GetStub().PutState(key, data)
}

// GetState returns state from world state.
// Unmarshalls the JSON into passed state.
// Key is the split key value used in Add/Update joined using a colon
func (sl *StateList) GetState(key string, state StateInterface) error {
	ledgerKey, _ := sl.Ctx.GetStub().CreateCompositeKey(sl.Name, SplitKey(key))
	data, err := sl.Ctx.GetStub().GetState(ledgerKey)

	if err != nil {
		return err
	} else if data == nil {
		return fmt.Errorf("No state found for %s", key)
	}

	return sl.Deserialize(data, state)
}

// GetAllStates returns all states from world state.
func (sl *StateList) GetAllStates() (shim.StateQueryIteratorInterface, error) {
	// As composite keys have been used, getStateByRange method won't work because of the \u0000 delimiter hyperledger uses.
	// Therefore for this implementation GetStateByPartialCompositeKey has been used and for each key "MedStore" string
	// has been attached for easier retrieval.
	resultsIterator, err := sl.Ctx.GetStub().GetStateByPartialCompositeKey(sl.Name, []string{"MedStore"})
	if err != nil {
		return nil, err
	}

	return resultsIterator, nil
}

// UpdateState puts state into world state. Same as AddState but
// separate as semantically different
func (sl *StateList) UpdateState(state StateInterface) error {
	return sl.AddState(state)
}
