package medicalsupply

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-tpm/tpm2"
	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/ledger-api"
)

// createTPMledgerKey - Creates a key for the TPM Authentication.
func createTPMledgerKey(holder string) string {
	return ledgerapi.MakeKey("TPMAUTH", holder)
}

type tpmAuthAlias TPMAuth
type jsonTPMAuth struct {
	*tpmAuthAlias
	Class string `json:"class"`
	Key   string `json:"key"`
}

type TPMAuth struct {
	Holder string `json:"holder"`
	TPMKey string `json:"tpmkey"`
	class  string `metadata:"class"`
	key    string `metadata:"key"`
}

//-------------------------------------------------------//

// MarshalJSON - Special handler for managing JSON marshalling.
func (auth TPMAuth) MarshalJSON() ([]byte, error) {
	jauth := jsonTPMAuth{tpmAuthAlias: (*tpmAuthAlias)(&auth), Class: "org.medstore.tpmauth", Key: createTPMledgerKey(auth.Holder)}
	return json.Marshal(&jauth)
}

// UnmarshalJSON - Special handler for managing JSON marshalling.
func (auth *TPMAuth) UnmarshalJSON(data []byte) error {
	jauth := jsonTPMAuth{tpmAuthAlias: (*tpmAuthAlias)(auth)}

	err := json.Unmarshal(data, &jauth)
	if err != nil {
		return err
	}
	return nil
}

// GetSplitKey - Returns values which should be used to form key.
func (auth *TPMAuth) GetSplitKey() []string {
	return []string{"TPMAUTH", auth.Holder}
}

// Serialize - Formats the tpm authentication as JSON bytes.
func (auth *TPMAuth) Serialize() ([]byte, error) {
	return json.Marshal(auth)
}

// Deserialize - Formats the tpm authentication from JSON bytes.
func DeserializeTPM(bytes []byte, auth *TPMAuth) error {
	err := json.Unmarshal(bytes, auth)

	if err != nil {
		return fmt.Errorf("error deserializing tpm authentication. %s", err.Error())
	}

	return nil
}

//-------------------------------------------------------//

// tpmHash - Hashes string using TPM 2.0.
func tpmHash(input string) (string, error) {
	input = strings.ToLower(input)

	// Sudo chown krn /dev/tpm0
	rwc, err := tpm2.OpenTPM("/dev/tpmrm0")
	if err != nil {
		return "", fmt.Errorf("couldn't open the TPM file /dev/tpm0: %s", err)
	}

	// Convert input to bytes.
	dataToHash := []byte(input)

	// Hash the bytes input.
	hashDigest, _, hashErr := tpm2.Hash(rwc, tpm2.AlgSHA256, dataToHash, tpm2.HandleOwner)
	if hashErr != nil {
		return "", fmt.Errorf("hash failed unexpectedly: %s", err)
	}

	// Close TPM rwc
	err = rwc.Close()
	if err != nil {
		return "", fmt.Errorf("error occurred when closing rwc for TPM: %s", err)
	}
	return string(hashDigest[:]), nil
}

// TpmKey - Generates key using TPM 2.0.
// Based on Go-TPM examples.
func tpmKey() (string, error) {
	pcrSelection7 := tpm2.PCRSelection{Hash: tpm2.AlgSHA1, PCRs: []int{7}} // 7 for secure boot

	rwc, err := tpm2.OpenTPM("/dev/tpmrm0")
	if err != nil {
		return "", fmt.Errorf("couldn't open the TPM file /dev/tpm0: %s", err)
	}

	// Generate random 16 bytes.
	randByte, err := tpm2.GetRandom(rwc, 16)
	if err != nil {
		return "", fmt.Errorf("generating random bytes for key failed: %s", err)
	}
	randompassword := string(randByte)

	// Generate primary key using the random bytes.
	parentHandle, _, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, pcrSelection7, "", randompassword, tpm2.Public{
		Type:       tpm2.AlgRSA,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagRestricted | tpm2.FlagDecrypt | tpm2.FlagUserWithAuth | tpm2.FlagFixedParent | tpm2.FlagFixedTPM | tpm2.FlagSensitiveDataOrigin,
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
			KeyBits: 2048,
		},
	})
	if err != nil {
		return "", fmt.Errorf("CreatePrimary failed: %s", err)
	}
	// Remove parenthandle from TPM to avoid out-of-memory problems.
	defer tpm2.FlushContext(rwc, parentHandle)

	// Create key using the handle of the primary key.
	privateBlob, publicBlob, _, _, _, err := tpm2.CreateKey(rwc, parentHandle, pcrSelection7, randompassword, randompassword, tpm2.Public{
		Type:       tpm2.AlgSymCipher,
		NameAlg:    tpm2.AlgSHA256,
		Attributes: tpm2.FlagDecrypt | tpm2.FlagSign | tpm2.FlagUserWithAuth | tpm2.FlagFixedParent | tpm2.FlagFixedTPM | tpm2.FlagSensitiveDataOrigin,
		SymCipherParameters: &tpm2.SymCipherParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("createKey failed: %s", err)
	}

	// Load the key created into an object and return its handle.
	createdkey, _, err := tpm2.Load(rwc, parentHandle, randompassword, publicBlob, privateBlob)
	if err != nil {
		return "", fmt.Errorf("loading key failed %s", err)
	}
	// Remove key from TPM to avoid out-of-memory problems.
	defer tpm2.FlushContext(rwc, createdkey)

	datablob := bytes.Repeat([]byte("a"), 8) // 8 bytes can be changed to increase the size of the key
	initvect := make([]byte, 16)             //16 byte long array

	// Create symmetric encryption key
	encrypted, err := tpm2.EncryptSymmetric(rwc, randompassword, createdkey, initvect, datablob)
	if err != nil {
		return "", fmt.Errorf("creating encryption key failed %s", err)
	}

	err = rwc.Close()
	if err != nil {
		return "", fmt.Errorf("error occurred when closing rwc for TPM: %s", err)
	}

	return string(encrypted[:]), nil
}
