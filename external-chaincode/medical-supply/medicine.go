package medicalsupply

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/ledger-api"
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
	CheckSum   string `json:"checkSum"`
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
	jms := jsonMedicalSupply{medicalSupplyAlias: (*medicalSupplyAlias)(&ms), State: ms.state, Class: "org.medstore.medicalsupply", Key: CreateMedicalKey(ms.MedName, ms.MedNumber)}

	return json.Marshal(&jms)
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

// InitialiseChecksum - Initialise the checksum value of the MedicalSupply using tpm hashing.
func (ms *MedicalSupply) InitialiseChecksum() error {
	checkSumStr := fmt.Sprintf(ms.MedName, ms.MedNumber, ms.Disease, ms.Expiration, ms.Price, ms.Holder)
	checksum, tpmError := tpmHash(checkSumStr)

	if tpmError != nil {
		return fmt.Errorf("tpm error occurred: %s", tpmError)
	}
	ms.CheckSum = checksum
	return nil
}

// VerifyChecksum - Returns true if the checksum stored on the Medicine object still is the same as after recalculating the checksum.
func (ms *MedicalSupply) VerifyChecksum() error {
	checkSumStr := fmt.Sprintf(ms.MedName, ms.MedNumber, ms.Disease, ms.Expiration, ms.Price)
	checksum, tpmError := tpmHash(checkSumStr)

	if tpmError != nil {
		return fmt.Errorf("tpm error occurred: %s", tpmError)
	}
	if ms.CheckSum != checksum {
		return fmt.Errorf("medical-supply is not valid to transaction for due to failed checksum")
	}
	return nil
}

// Serialize - Formats the medical supply as JSON bytes.
func (ms *MedicalSupply) Serialize() ([]byte, error) {
	return json.Marshal(ms)
}

// Deserialize - Formats the commercial paper from JSON bytes.
func DeserializeJSON(bytes []byte, ms *MedicalSupply) error {
	err := json.Unmarshal(bytes, ms)

	if err != nil {
		return fmt.Errorf("error deserializing medical supply. %s", err.Error())
	}

	return nil
}
