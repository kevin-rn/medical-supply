package medicalsupply

import (
	"testing"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/regulators/chaincode/ledger-api"
	"github.com/stretchr/testify/assert"
)

func TestCreateTPMledgerKey(t *testing.T) {
	assert.Equal(t, ledgerapi.MakeKey("TPMAUTH", "hashedusername"), createTPMledgerKey("hashedusername"), "should return key comprised of passed values.")
}

func TestGetTPMSplitKey(t *testing.T) {
	auth := new(TPMAuth)
	auth.Holder = "hashedusername"

	assert.Equal(t, []string{"TPMAUTH", "hashedusername"}, auth.GetSplitKey(),
		"should return TPMAUTH and hashedusername as split key.")
}

func TestSerializeTPM(t *testing.T) {
	auth := new(TPMAuth)
	auth.Holder = "hashedusername"
	auth.TPMKey = "hashedkey"
	correctJson := `{"holder":"hashedusername","tpmkey":"hashedkey","class":"org.medstore.tpmauth","key":"TPMAUTH:hashedusername"}`

	bytes, err := auth.Serialize()
	assert.Nil(t, err, "should not error on serialize")
	assert.Equal(t, correctJson, string(bytes), "should return JSON formatted value")
}

func TestDeserializeTPM(t *testing.T) {
	var auth *TPMAuth
	var err error

	auth = new(TPMAuth)
	correctJson := `{"holder":"hashedusername","tpmkey":"hashedkey","class":"org.medstore.tpmauth","key":"TPMAUTH:hashedusername"}`
	err = DeserializeTPM([]byte(correctJson), auth)
	assert.Nil(t, err, "should not return error for deserialize")

	expectedAuth := new(TPMAuth)
	expectedAuth.Holder = "hashedusername"
	expectedAuth.TPMKey = "hashedkey"
	assert.Equal(t, expectedAuth, auth, "should create expected tpm authentication")

	incorrectJson := `{"holder":10, "key":"hashedkey", "class":"org.medstore.tpmauth","key":"TPMAUTH:hashedusername"}`
	auth = new(TPMAuth)
	err = DeserializeTPM([]byte(incorrectJson), auth)
	assert.EqualError(t, err, "Error deserializing tpm authentication. json: cannot unmarshal number into Go struct field jsonTPMAuth.holder of type string", "should return error for bad data")
}
