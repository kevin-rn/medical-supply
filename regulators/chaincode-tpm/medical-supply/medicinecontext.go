package medicalsupply

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to describe
// the minimum required functions for a transaction context.
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetMedicineList() ListInterface
}

// TransactionContext implementation of TransactionContextInterface
// for use with contract.
type TransactionContext struct {
	contractapi.TransactionContext
	medicineList *list
}

// GetMedicineList return medicine list (represents the ledger).
func (tc *TransactionContext) GetMedicineList() ListInterface {
	if tc.medicineList == nil {
		tc.medicineList = newList(tc)
	}

	return tc.medicineList
}
