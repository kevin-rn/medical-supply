package medicalsupply

import (
	"testing"

	ledgerapi "github.com/hyperledger/fabric-samples/medical-supply/customers/chaincode/ledger-api"
	"github.com/stretchr/testify/assert"
)

func TestGetMedicineList(t *testing.T) {
	var tc *TransactionContext
	var expectedMedicineList *list

	tc = new(TransactionContext)
	expectedMedicineList = newList(tc)
	actualList := tc.GetMedicineList().(*list)
	assert.Equal(t, expectedMedicineList.statelist.(*ledgerapi.StateList).Name, actualList.statelist.(*ledgerapi.StateList).Name, "should configure medicine list when one not already configured")

	tc = new(TransactionContext)
	expectedMedicineList = new(list)
	expectedStateList := new(ledgerapi.StateList)
	expectedStateList.Ctx = tc
	expectedStateList.Name = "existing medicine list"
	expectedMedicineList.statelist = expectedStateList
	tc.medicineList = expectedMedicineList
	assert.Equal(t, expectedMedicineList, tc.GetMedicineList(), "should return set medicine list when already set")
}
