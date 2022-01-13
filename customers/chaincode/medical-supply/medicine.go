package medicalsupply

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/ledger-api"
)

type State uint

const (
	// AVAILABLE state for when medicine has been issued.
	AVAILABLE State = iota + 1
	// REQUESTED state for when medicine has been requested.
	REQUESTED
	// SEND state for when a medicine is send.
	SEND
)

// String - Changes state enum to string.
func (state State) String() string {
	names := []string{"AVAILABLE", "REQUESTED", "SEND"}

	if state < AVAILABLE || state > SEND {
		return "UNKNOWN"
	}
	return names[state-1]
}

// CreateMedicalKey - Creates a key for the medical supply (e.g. MedStore:Aspirin:00001).
func CreateMedicalKey(medName string, medNumber string) string {
	return ledgerapi.MakeKey("MedStore", medName, medNumber)
}

// Used for managing the fact state is private but still used it in the world state.
type medicalSupplyAlias MedicalSupply
type jsonMedicalSupply struct {
	*medicalSupplyAlias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// MedicalSupply - Defines a medicine.
type MedicalSupply struct {
	MedName    string `json:"medName"`
	MedNumber  string `json:"medNumber"`
	Disease    string `json:"disease"`
	Expiration string `json:"expiration"`
	Price      string `json:"price"`
	Holder     string `json:"holder"`
	state      State  `metadata:"currentState"`
	class      string `metadata:"class"`
	key        string `metadata:"key"`
}

//-------------------------------------------------------//

// MarshalJSON - Special handler for managing JSON marshalling.
func (ms MedicalSupply) MarshalJSON() ([]byte, error) {
	jcp := jsonMedicalSupply{medicalSupplyAlias: (*medicalSupplyAlias)(&ms), State: ms.state, Class: "org.medstore.medicalsupply", Key: ledgerapi.MakeKey("MedStore", ms.MedName, ms.MedNumber)}

	return json.Marshal(&jcp)
}

// UnmarshalJSON - Special handler for managing JSON marshalling.
func (ms *MedicalSupply) UnmarshalJSON(data []byte) error {
	jms := jsonMedicalSupply{medicalSupplyAlias: (*medicalSupplyAlias)(ms)}

	err := json.Unmarshal(data, &jms)
	if err != nil {
		return err
	}

	ms.state = jms.State
	return nil
}

//-------------------------------------------------------//

// GetState - Returns the state.
func (ms *MedicalSupply) GetState() State {
	return ms.state
}

// SetAvailable - Returns the state to AVAILABLE.
func (ms *MedicalSupply) SetAvailable() {
	ms.state = AVAILABLE
}

// SetRequested - Returns the state to REQUESTED.
func (ms *MedicalSupply) SetRequested() {
	ms.state = REQUESTED
}

// SetSend - Returns the state to SEND.
func (ms *MedicalSupply) SetSend() {
	ms.state = SEND
}

// IsAvailable - Returns true if state is AVAILABLE.
func (ms *MedicalSupply) IsAvailable() bool {
	return ms.state == AVAILABLE
}

// IsRequested - Returns true if state is REQUESTED.
func (ms *MedicalSupply) IsRequested() bool {
	return ms.state == REQUESTED
}

// IsSend - Returns true if state is SEND.
func (ms *MedicalSupply) IsSend() bool {
	return ms.state == SEND
}

//-------------------------------------------------------//

// GetSplitKey - Returns values which should be used to form key.
func (ms *MedicalSupply) GetSplitKey() []string {
	return []string{"MedStore", ms.MedName, ms.MedNumber}
}

// Serialize - Formats the medical supply as JSON bytes.
func (ms *MedicalSupply) Serialize() ([]byte, error) {
	return json.Marshal(ms)
}

// Deserialize - Formats the commercial paper from JSON bytes.
func Deserialize(bytes []byte, ms *MedicalSupply) error {
	err := json.Unmarshal(bytes, ms)

	if err != nil {
		return fmt.Errorf("Error deserializing medical supply. %s", err.Error())
	}

	return nil
}
