package medicalsupply

import (
	"testing"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/ledger-api"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, "AVAILABLE", AVAILABLE.String(), "should return string for available.")
	assert.Equal(t, "REQUESTED", REQUESTED.String(), "should return string for requested.")
	assert.Equal(t, "SEND", SEND.String(), "should return string for send.")
	assert.Equal(t, "UNKNOWN", State(SEND+1).String(), "should return unknown when not one of constants.")
}

func TestCreateMedicalKey(t *testing.T) {
	assert.Equal(t, ledgerapi.MakeKey("MedStore", "medicinename", "medicinenumber"), CreateMedicalKey("medicinename", "medicinenumber"), "should return key comprised of passed values.")
}

func TestGetState(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.SetAvailable()

	assert.Equal(t, AVAILABLE, medicine.GetState(), "should return set state.")
}

func TestSetAvailable(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.SetAvailable()
	assert.Equal(t, AVAILABLE, medicine.GetState(), "should set state to available.")
}

func TestSetRequested(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.SetRequested()
	assert.Equal(t, REQUESTED, medicine.GetState(), "should set state to requested.")
}

func TestSetSend(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.SetSend()
	assert.Equal(t, SEND, medicine.GetState(), "should set state to send.")
}

func TestIsIssued(t *testing.T) {
	medicine := new(MedicalSupply)

	medicine.SetAvailable()
	assert.True(t, medicine.IsAvailable(), "should be true when status set to available.")

	medicine.SetRequested()
	assert.False(t, medicine.IsAvailable(), "should be false when status not set to available.")
}

func TestIsRequested(t *testing.T) {
	medicine := new(MedicalSupply)

	medicine.SetRequested()
	assert.True(t, medicine.IsRequested(), "should be true when status set to requested.")

	medicine.SetSend()
	assert.False(t, medicine.IsRequested(), "should be false when status not set to requested.")
}

func TestIsSend(t *testing.T) {
	medicine := new(MedicalSupply)

	medicine.SetSend()
	assert.True(t, medicine.IsSend(), "should be true when status set to send.")

	medicine.SetRequested()
	assert.False(t, medicine.IsSend(), "should be false when status not set to send.")
}

func TestGetSplitKey(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.MedName = "medicinename"
	medicine.MedNumber = "medicinenumber"

	assert.Equal(t, []string{"MedStore", "medicinename", "medicinenumber"}, medicine.GetSplitKey(),
		"should return medicine name and number as split key.")
}

func TestSerialize(t *testing.T) {
	medicine := new(MedicalSupply)
	medicine.MedName = "aspirin"
	medicine.MedNumber = "00001"
	medicine.Disease = "pain"
	medicine.Expiration = "2022.02.22"
	medicine.Price = "$10"
	medicine.Holder = "alice"
	medicine.SetAvailable()

	correctJson := `{"checkSum":"","medName":"aspirin","medNumber":"00001","disease":"pain","expiration":"2022.02.22","price":"$10","holder":"alice","currentState":1,"class":"org.medstore.medicalsupply","key":"MedStore:aspirin:00001"}`

	bytes, err := medicine.Serialize()
	assert.Nil(t, err, "should not error on serialize")
	assert.Equal(t, correctJson, string(bytes), "should return JSON formatted value")
}

func TestDeserialize(t *testing.T) {
	var medicine *MedicalSupply
	var err error

	medicine = new(MedicalSupply)
	correctJson := `{"medName":"aspirin","medNumber":"00001","disease":"pain","expiration":"2022.02.22","price":"$10","holder":"alice","currentState":1,"class":"org.medstore.medicalsupply","key":"MedStore:aspirin:00001"}`
	err = DeserializeJSON([]byte(correctJson), medicine)
	assert.Nil(t, err, "should not return error for deserialize")

	expectedMedicine := new(MedicalSupply)
	expectedMedicine.MedName = "aspirin"
	expectedMedicine.MedNumber = "00001"
	expectedMedicine.Disease = "pain"
	expectedMedicine.Expiration = "2022.02.22"
	expectedMedicine.Price = "$10"
	expectedMedicine.Holder = "alice"
	expectedMedicine.SetAvailable()
	assert.Equal(t, expectedMedicine, medicine, "should create expected medical supply")

	incorrectJson := `{"medName":"aspirin","medNumber":"00001","disease": 404,"expiration":"2022.02.22","price":"$10","holder":"alice","currentState":1,"class":"org.medstore.medicalsupply","key":"MedStore:aspirin:00001"}`
	medicine = new(MedicalSupply)
	err = DeserializeJSON([]byte(incorrectJson), medicine)
	assert.EqualError(t, err, "Error deserializing medical supply. json: cannot unmarshal number into Go struct field jsonMedicalSupply.disease of type string", "should return error for bad data")
}
