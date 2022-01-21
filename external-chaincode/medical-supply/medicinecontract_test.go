package medicalsupply

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/hyperledger/fabric-chaincode-go/shim"
// 	"github.com/hyperledger/fabric-samples/medical-supply/external-chaincode/medical-supply/mocks"
// )

// //go:generate counterfeiter -o mocks/transactioncontext.go -fake-name TransactionContext . transactionContext
// type transactionContext interface {
// 	TransactionContextInterface
// }

// //go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
// type chaincodeStub interface {
// 	shim.ChaincodeStubInterface
// }

// //go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
// type stateQueryIterator interface {
// 	shim.StateQueryIteratorInterface
// }

// func TestInitLedger(t *testing.T) {
// 	chaincodeStub := &mocks.ChaincodeStub{}
// 	transactionContext := &mocks.TransactionContext{}
// 	transactionContext.GetStubReturns(chaincodeStub)

// 	// chaincode := Contract{}
// 	// err := chaincode.InitLedger(transactionContext)
// 	// require.NoError(t, err)

// 	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
// 	// err = chaincode.InitLedger(transactionContext)
// 	// require.EqualError(t, err, "failed to put to world state. failed inserting key")
// }

// func TestCreateAsset(t *testing.T) {
// 	chaincodeStub := &mocks.ChaincodeStub{}
// 	transactionContext := &mocks.TransactionContext{}
// 	transactionContext.GetStubReturns(chaincodeStub)

// 	assetTransfer := Contract{}
// 	err := assetTransfer.Issue(transactionContext, "", "", 0, "", 0)
// 	require.NoError(t, err)

// 	chaincodeStub.GetStateReturns([]byte{}, nil)
// 	err = assetTransfer.Issue(transactionContext, "asset1", "", 0, "", 0)
// 	require.EqualError(t, err, "the asset asset1 already exists")

// 	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
// 	err = assetTransfer.Issue(transactionContext, "asset1", "", 0, "", 0)
// 	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
// }
