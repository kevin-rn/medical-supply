package medicalsupply

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/stakeholders/customers/chaincode-go/ledger-api"
)

type State uint

const (
	// REQUESTED state for when medicine has been requested
	REQUESTED State = iota + 1
	// SENT state for when a medicine is sent
	SENT
)

func (state State) String() string {
	names := []string{"REQUESTED", "SENT"}

	if state < REQUESTED || state > SENT {
		return "UNKNOWN"
	}
	return names[state-1]
}

// CreateMedicalKey - Creates a key for the medical supply (e.g. MedStoreAspirin0000)
func CreateMedicalKey(medName string, medNumber string) string {
	return ledgerapi.MakeKey("MedStore", medName, medNumber)
}

// Used for managing the fact status is private but want it in the world state.
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
	jcp := jsonMedicalSupply{medicalSupplyAlias: (*medicalSupplyAlias)(&ms), State: ms.state, Class: "org.medstore.medicalsupply", Key: ledgerapi.MakeKey(ms.MedNumber)}

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

// GetState returns the state
func (ms *MedicalSupply) GetState() State {
	return ms.state
}

// SetRequested returns the state to requested
func (ms *MedicalSupply) SetRequested() {
	ms.state = REQUESTED
}

// SetSent returns the state to set
func (ms *MedicalSupply) SetSent() {
	ms.state = SENT
}

// IsRequested returns true if state is issued
func (ms *MedicalSupply) IsRequested() bool {
	return ms.state == REQUESTED
}

// IsSent returns true if state is sent
func (ms *MedicalSupply) IsSent() bool {
	return ms.state == SENT
}

//-------------------------------------------------------//

// GetSplitKey returns values which should be used to form key
func (ms *MedicalSupply) GetSplitKey() []string {
	return []string{ms.MedName, ms.MedNumber}
}

//-------------------------------------------------------//

// Serialize formats the medical supply as JSON bytes
func (ms *MedicalSupply) Serialize() ([]byte, error) {
	return json.Marshal(ms)
}

// Deserialize formats the commercial paper from JSON bytes
func Deserialize(bytes []byte, ms *MedicalSupply) error {
	err := json.Unmarshal(bytes, ms)

	if err != nil {
		return fmt.Errorf("Error deserializing medical supply. %s", err.Error())
	}

	return nil
}
