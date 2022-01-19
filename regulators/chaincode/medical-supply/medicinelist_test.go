package medicalsupply

import (
	"testing"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/regulators/chaincode/ledger-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStateList struct {
	mock.Mock
}

func TestNewStateList(t *testing.T) {
	ctx := new(TransactionContext)
	list := newList(ctx)
	stateList, ok := list.statelist.(*ledgerapi.StateList)

	assert.True(t, ok, "should make statelist of type ledgerapi.StateList")
	assert.Equal(t, ctx, stateList.Ctx, "should set the context to passed context")
	assert.Equal(t, "org.medstore.medicalsupplylist", stateList.Name, "should set the name for the list")

	expectedErr := Deserialize([]byte("bad json"), new(MedicalSupply))
	err := stateList.Deserialize([]byte("bad json"), new(MedicalSupply))
	assert.EqualError(t, err, expectedErr.Error(), "should call Deserialize when stateList.Deserialize called")
}
