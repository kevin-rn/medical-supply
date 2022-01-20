package medicalsupply

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

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
	Holder  string `json:"holder"`
	TPMKey  string `json:"tpmkey"`
	IsAdmin bool   `json:"isadmin"`
	class   string `metadata:"class"`
	key     string `metadata:"key"`
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
		return fmt.Errorf("Error deserializing tpm authentication. %s", err.Error())
	}

	return nil
}

//-------------------------------------------------------//

var (
	handleNames = map[string][]tpm2.HandleType{
		"all":       []tpm2.HandleType{tpm2.HandleTypeLoadedSession, tpm2.HandleTypeSavedSession, tpm2.HandleTypeTransient},
		"loaded":    []tpm2.HandleType{tpm2.HandleTypeLoadedSession},
		"saved":     []tpm2.HandleType{tpm2.HandleTypeSavedSession},
		"transient": []tpm2.HandleType{tpm2.HandleTypeTransient},
	}

	tpmPath = flag.String("tpm-path", "/dev/tpm0", "Path to the TPM device (character device or a Unix socket).")

	defaultEKTemplate = tpm2.Public{
		Type:    tpm2.AlgRSA,
		NameAlg: tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
			tpm2.FlagAdminWithPolicy | tpm2.FlagRestricted | tpm2.FlagDecrypt,
		AuthPolicy: []byte{
			0x83, 0x71, 0x97, 0x67, 0x44, 0x84,
			0xB3, 0xF8, 0x1A, 0x90, 0xCC, 0x8D,
			0x46, 0xA5, 0xD7, 0x24, 0xFD, 0x52,
			0xD7, 0x6E, 0x06, 0x52, 0x0B, 0x64,
			0xF2, 0xA1, 0xDA, 0x1B, 0x33, 0x14,
			0x69, 0xAA,
		},
		RSAParameters: &tpm2.RSAParams{
			Symmetric: &tpm2.SymScheme{
				Alg:     tpm2.AlgAES,
				KeyBits: 128,
				Mode:    tpm2.AlgCFB,
			},
			KeyBits:    2048,
			ModulusRaw: make([]byte, 256),
		},
	}

	// https://github.com/google/go-tpm/blob/master/tpm2/constants.go#L152
	defaultKeyParams = tpm2.Public{
		Type:    tpm2.AlgRSA,
		NameAlg: tpm2.AlgSHA256,
		Attributes: tpm2.FlagSign | tpm2.FlagRestricted | tpm2.FlagFixedTPM |
			tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin | tpm2.FlagUserWithAuth,
		AuthPolicy: []byte{},
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
			KeyBits: 2048,
		},
	}

	unrestrictedKeyParams = tpm2.Public{
		Type:    tpm2.AlgRSA,
		NameAlg: tpm2.AlgSHA256,
		Attributes: tpm2.FlagFixedTPM | tpm2.FlagFixedParent | tpm2.FlagSensitiveDataOrigin |
			tpm2.FlagUserWithAuth | tpm2.FlagSign,
		AuthPolicy: []byte{},
		RSAParameters: &tpm2.RSAParams{
			Sign: &tpm2.SigScheme{
				Alg:  tpm2.AlgRSASSA,
				Hash: tpm2.AlgSHA256,
			},
			KeyBits: 2048,
		},
	}
)

func tpmHash(input string) (string, error) {
	// Sudo chown kevin /dev/tpm0
	flag.Parse()
	rwc, err := tpm2.OpenTPM(*tpmPath)
	if err != nil {
		return "", err
	}

	err = rwc.Close()
	if err != nil {
		return "", err
	}

	dataToHash := []byte(input)
	hashDigest, _, hashErr := tpm2.Hash(rwc, tpm2.AlgSHA256, dataToHash, tpm2.HandleOwner)
	if hashErr != nil {
		log.Fatalf("Hash failed unexpectedly: %v", err)
	}

	return string(hashDigest[:]), nil
}
