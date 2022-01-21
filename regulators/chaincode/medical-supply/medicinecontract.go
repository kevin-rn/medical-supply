package medicalsupply

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) Instantiate() {
	fmt.Println("- Contract Instantiated -")
}

// TPMKeyGen - Helper function for generating tpm key to be used for authentication.
// Only returns key at first creation as calling this function repeatedly would otherwise be exploitable.
func (c *Contract) TPMKeyGen(ctx TransactionContextInterface, user string) (string, error) {
	bool := ctx.GetMedicineList().ExistsTPMAuth(user)
	if bool {

		tpmkey, err := tpmKey()
		if err != nil {
			tpmkey = "password" + user //Temporary solution
			// return "", fmt.Errorf("could not generate tpm key: %s", err)
		}
		// user, err = tpmHash(user)
		// if err != nil {
		// 	return "", fmt.Errorf("could hash user name: %s", err)
		// }

		// Create MedicalSupply object.
		tpmAuth := TPMAuth{Holder: user, TPMKey: tpmkey}
		err = ctx.GetMedicineList().AddTPMAuth(&tpmAuth)
		if err != nil {
			return "", err
		}
		return tpmAuth.TPMKey, nil
	}
	return "", fmt.Errorf("user %s has already created a TPM authentication", user)
}

// tpmCheck - Helper function for verifying authentication
func (c *Contract) tpmCheck(ctx TransactionContextInterface, user string, tpmkey string) error {
	check, err := ctx.GetMedicineList().VerifyTPMAuth(user, tpmkey)
	if err != nil {
		return fmt.Errorf("user has not authenticated yet. Please invoke TPMKeyGen first: %s", err)
	}
	if check {
		return fmt.Errorf("provided tpm key does not match with registered authentication: %s", err)
	}
	return nil
}

// hasAuthority - Helper function for verifying the invoker organisation.
func (c *Contract) hasAuthority(ctx TransactionContextInterface, user string, tpmkey string) error {
	// Check if user is authenticated
	err := c.tpmCheck(ctx, user, tpmkey)
	if err != nil {
		return err
	}

	ciMsp, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if ciMsp != "Org2MSP" {
		return fmt.Errorf("user from organisation %s, does not have acces to this function", ciMsp)
	}
	return nil
}

// InitLedger - Adds a base set of medicine (MedicalSupply) to the ledger. [Regulators]
func (c *Contract) InitLedger(ctx TransactionContextInterface, user string, tpmkey string) error {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return err
	}

	// Create array of MedicalSupply objects.
	medicines := []MedicalSupply{
		{MedName: "aspirin", MedNumber: "00001", Disease: "pain management", Expiration: "2022.05.09", Price: "$10", Holder: "MedStore"},
		{MedName: "vicodin", MedNumber: "00002", Disease: "pain management", Expiration: "2022.07.01", Price: "$14", Holder: "MedStore"},
		{MedName: "synthroid", MedNumber: "00003", Disease: "thyroid deficiency", Expiration: "2021.12.03", Price: "$11", Holder: "MedStore"},
		{MedName: "delasone", MedNumber: "00004", Disease: "arthritis", Expiration: "2022.09.12", Price: "$5", Holder: "MedStore"},
		{MedName: "amoxil", MedNumber: "00005", Disease: "bacterial infections", Expiration: "2022.07.08", Price: "$9", Holder: "MedStore"},
		{MedName: "neurontin", MedNumber: "00006", Disease: "seizures", Expiration: "2022.03.25", Price: "$13", Holder: "MedStore"},
		{MedName: "zestril", MedNumber: "00007", Disease: "blood pressure", Expiration: "2022.03.11", Price: "$7", Holder: "MedStore"},
		{MedName: "lipitor", MedNumber: "00008", Disease: "high cholesterol", Expiration: "2022.01.06", Price: "$12", Holder: "MedStore"},
		{MedName: "glucophage", MedNumber: "00009", Disease: "type 2 diabetes", Expiration: "2022.04.24", Price: "$8", Holder: "MedStore"},
		{MedName: "zofran", MedNumber: "00010", Disease: "fever", Expiration: "2022.02.04", Price: "$13", Holder: "MedStore"},
		{MedName: "ibuprofen", MedNumber: "00011", Disease: "fever", Expiration: "2022.02.28", Price: "$12", Holder: "MedStore"},
	}

	// For each medicine, set it's state to Available, calculate the checksum and update the ledger
	for _, med := range medicines {
		med.SetAvailable()
		// TODO: checksum
		err := ctx.GetMedicineList().UpdateMedicine(&med)

		if err != nil {
			return fmt.Errorf("failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// Issue - Function for handling issued medicine [Regulators]
func (c *Contract) Issue(ctx TransactionContextInterface, medname string, mednumber string,
	disease string, expiration string, price string, user string, tpmkey string) (*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Create MedicalSupply object.
	medicine := MedicalSupply{
		MedName:    strings.ToLower(medname),
		MedNumber:  mednumber,
		Disease:    strings.ToLower(disease),
		Expiration: expiration,
		Price:      price,
		Holder:     "MedStore",
	}

	// Calculate the checksum by using the hashfunction of the TPM.
	err = medicine.InitialiseChecksum()
	if err != nil {
		return nil, fmt.Errorf("could not issue new MedicalSupply. %s", err)
	}

	// Set state to AVAILABLE.
	medicine.SetAvailable()

	// Add the medicine to the ledger.
	err = ctx.GetMedicineList().AddMedicine(&medicine)
	if err != nil {
		return nil, fmt.Errorf("could not add medicine to the ledger: %s", err)
	}

	return &medicine, nil
}

// Delete - Function for handling medicine removal. [Regulators]
func (c *Contract) Delete(ctx TransactionContextInterface, medName string, medNumber string, user string, tpmkey string) error {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	if medicine == nil {
		return fmt.Errorf("medicine does not exist, can't delete from ledger")
	}
	return ctx.GetMedicineList().DeleteMedicine(medName, medNumber)
}

// Request - Function for handling requested medicine. [Customers]
func (c *Contract) Request(ctx TransactionContextInterface, medName string, medNumber string, user string, tpmkey string) (*MedicalSupply, error) {
	// Hashes user string
	// user, err := tpmHash(user)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot hash user string: %s", err)
	// }

	// Checks authentication
	err := c.tpmCheck(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// Verify that the current holder is MedStore, if that is not the case than the medicine has already been transferred to a different holder.
	if medicine.Holder != "MedStore" {
		return nil, fmt.Errorf("medicine %s:%s has already been bought", medName, medNumber)
	}

	// Verify that the current state is AVAILABLE, if so set to REQUESTED.
	if medicine.IsAvailable() {
		medicine.SetRequested()
	} else {
		return nil, fmt.Errorf("medicine %s:%s is currently not available at MedStore", medName, medNumber)
	}

	// Verify that change to REQUESTED state has succeeded.
	if !medicine.IsRequested() {
		return nil, fmt.Errorf("medicine %s:%s is not requested. current state = %s", medName, medNumber, medicine.GetState())
	}

	// Update medicine holder to be the customer instead of MedStore.
	medicine.Holder = user
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}

// CancelRequest - Function for handling cancelled requested medicine. [Customers]
func (c *Contract) CancelRequest(ctx TransactionContextInterface, medName string, medNumber string, user string, tpmkey string) (*MedicalSupply, error) {
	// Hashes user string
	// user, err := tpmHash(user)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot hash user string: %s", err)
	// }

	// Checks authentication
	err := c.tpmCheck(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// Check if medicine state is REQUESTED, if so set it to AVAILABLE and reset to holder to be MedStore.
	if medicine.IsRequested() && medicine.Holder == user {
		medicine.SetAvailable()
		medicine.Holder = "MedStore"
	} else {
		return nil, fmt.Errorf("cannot cancel because medicine has not been requested")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}

// SearchMedicineByName - Function for getting information on available medicine given the medicine name. [Customers]
func (c *Contract) SearchMedicineByName(ctx TransactionContextInterface, medName string) ([]*MedicalSupply, error) {
	// Retrieve the medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicineByName(medName)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}
	// Loop through the list and check for AVAILABLE state.
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.IsAvailable() {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// CheckHistory - Function for getting an overview of all Medicine. [Regulators]
func (c *Contract) CheckHistory(ctx TransactionContextInterface, user string, tpmkey string) ([]*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve query any medicine from ledger: %s", err)
	}
	return medicinelist, nil
}

// CheckAvailableMedicine - Function for getting an overview of all available medicine. [Customers]
func (c *Contract) CheckAvailableMedicine(ctx TransactionContextInterface) ([]*MedicalSupply, error) {
	// Get all medicine from the ledger (There is currently no efficienter way to retrieve assets from the Ledger for certain fields).
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, fmt.Errorf("could not query any medicine from ledger: %s", err)
	}

	// Loop through the list and check for AVAILABLE state.
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.IsAvailable() {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// CheckRequestedMedicine - Function for getting an overview of all requested medicine. [Regulators]
func (c *Contract) CheckRequestedMedicine(ctx TransactionContextInterface, user string, tpmkey string) ([]*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, fmt.Errorf("could not query any medicine from ledger: %s", err)
	}

	// Loop through the list and check for REQUESTED state.
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.IsRequested() {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// CheckUserHistory - Function for getting an overview of all medicine an user has ordered. [Customers]
func (c *Contract) CheckUserHistory(ctx TransactionContextInterface, user string, tpmkey string) ([]*MedicalSupply, error) {
	// Hashes user string
	// user, err := tpmHash(user)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot hash user string: %s", err)
	// }

	// Checks authentication
	err := c.tpmCheck(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Get all medicine from the ledger.
	medicinelist, err := ctx.GetMedicineList().GetAllMedicine()
	if err != nil {
		return nil, fmt.Errorf("could not query any medicine from ledger: %s", err)
	}

	// Loop through the list and check for the user (holder).
	var resultlist []*MedicalSupply
	for _, med := range medicinelist {
		if med.Holder == user {
			resultlist = append(resultlist, med)
		}
	}
	return resultlist, nil
}

// ApproveRequest - Function for handling approving the medicine by changing its state to SEND. [Regulators]
func (c *Contract) ApproveRequest(ctx TransactionContextInterface, medName string, medNumber string, user string, tpmkey string) (*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// Check if medicine state is REQUESTED, if so set it to SEND.
	if medicine.IsRequested() {
		medicine.SetSend()
	} else {
		return nil, fmt.Errorf("cannot approve medicine that has not been requested")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}

// RejectRequest - Function for handling disapproving the medicine by changing its state back to AVAILABLE. [Regulators]
func (c *Contract) RejectRequest(ctx TransactionContextInterface, medName string, medNumber string, user string, tpmkey string) (*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// Check if medicine state is REQUESTED, if so set it to AVAILABLE and reset to holder to be MedStore.
	if medicine.IsRequested() {
		medicine.SetAvailable()
		medicine.Holder = "MedStore"
	} else {
		return nil, fmt.Errorf("cannot disapprove medicine that has not been requested")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}

// ChangeStatus - Function for changing the status of a medicine. [Regulators]
func (c *Contract) ChangeStatus(ctx TransactionContextInterface, medName string, medNumber string, status string, user string, tpmkey string) (*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// Match case on status and change it.
	switch strings.ToLower(status) {
	case "available":
		medicine.SetAvailable()
	case "requested":
		medicine.SetRequested()
	case "send":
		medicine.SetSend()
	default:
		return nil, fmt.Errorf("cannot change status to a non-possible state")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}

// ChangeHolder - Function for changing the holder of a medicine. [Regulators]
func (c *Contract) ChangeHolder(ctx TransactionContextInterface, medName string, medNumber string, customer string, user string, tpmkey string) (*MedicalSupply, error) {
	// Check acces rights
	err := c.hasAuthority(ctx, user, tpmkey)
	if err != nil {
		return nil, err
	}

	// Retrieve the medicine from the ledger.
	medicine, err := ctx.GetMedicineList().GetMedicine(medName, medNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve medicine from ledger: %s", err)
	}

	// customer, err = tpmHash(customer)
	// if err != nil {
	// 	return nil, err
	// }

	if len(customer) > 0 {
		medicine.Holder = customer
	} else {
		return nil, fmt.Errorf("can't change current holder to invalid username")
	}

	// Update medicine on the ledger
	err = ctx.GetMedicineList().UpdateMedicine(medicine)
	if err != nil {
		return nil, fmt.Errorf("could not update medicine on the ledger: %s", err)
	}

	return medicine, nil
}
