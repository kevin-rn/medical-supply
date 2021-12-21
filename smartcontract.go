package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing medicines
type SmartContract struct {
	contractapi.Contract
}

type Medicine struct {
	Name       string `json:"name"`
	Disease    string `json:"disease"`
	Expiration string `json:"expiration"`
	Price      string `json:"price"`
	Holder     string `json:"holder"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Medicine
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	medicines := []Medicine{
		Medicine{Name: "Aspirin", Disease: "Pain management", Expiration: "2022.05.09", Price: "$10", Holder: "MedStore"},
		Medicine{Name: "Vicodin", Disease: "Pain management", Expiration: "2022.07.01", Price: "$14", Holder: "MedStore"},
		Medicine{Name: "Synthroid", Disease: "Thyroid deficiency", Expiration: "2021.12.03", Price: "$11", Holder: "MedStore"},
		Medicine{Name: "Delasone", Disease: "Arthritis", Expiration: "2022.09.12", Price: "$5", Holder: "MedStore"},
		Medicine{Name: "Amoxil", Disease: "Bacterial infections", Expiration: "2022.07.08", Price: "$9", Holder: "MedStore"},
		Medicine{Name: "Neurontin", Disease: "Seizures", Expiration: "2022.03.25", Price: "$13", Holder: "MedStore"},
		Medicine{Name: "Zestril", Disease: "Blood pressure", Expiration: "2022.03.11", Price: "$7", Holder: "MedStore"},
		Medicine{Name: "Lipitor", Disease: "High cholesterol", Expiration: "2022.01.06", Price: "$12", Holder: "MedStore"},
		Medicine{Name: "Glucophage", Disease: "Type 2 diabetes", Expiration: "2022.04.24", Price: "$8", Holder: "MedStore"},
		Medicine{Name: "Zofran", Disease: "Nausea", Expiration: "2022.02.04", Price: "$13", Holder: "MedStore"},
		Medicine{Name: "Ibuprofen", Disease: "Fever", Expiration: "2022.02.28", Price: "$12", Holder: "MedStore"},
	}
	for i, med := range medicines {
		medAsBytes, _ := json.Marshal(med)
		err := ctx.GetStub().PutState("MED"+strconv.Itoa(i), medAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// createMed adds a new medicine to the world state with given details
func (s *SmartContract) createMed(ctx contractapi.TransactionContextInterface, medNumber string, name string, disease string, expiration string, price string, holder string) error {
	med := Medicine{
		Name:       name,
		Disease:    disease,
		Expiration: expiration,
		Price:      price,
		Holder:     holder,
	}

	medAsBytes, _ := json.Marshal(med)

	return ctx.GetStub().PutState(medNumber, medAsBytes)
}

// QueryMed returns the med stored in the world state with given id
func (s *SmartContract) QueryMed(ctx contractapi.TransactionContextInterface, medNumber string) (*Medicine, error) {
	medAsBytes, err := ctx.GetStub().GetState(medNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if medAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", medNumber)
	}

	med := new(Medicine)
	_ = json.Unmarshal(medAsBytes, med)

	return med, nil
}

// QueryAllMed returns all medicine found in world state
func (s *SmartContract) QueryAllMed(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		med := new(Medicine)
		_ = json.Unmarshal(queryResponse.Value, med)

		queryResult := QueryResult{Key: queryResponse.Key, Record: med}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeMedHolder updates the holder field of medicine with given id in world state
func (s *SmartContract) ChangeMedHolder(ctx contractapi.TransactionContextInterface, medNumber string, newHolder string) error {
	med, err := s.QueryMed(ctx, medNumber)

	if err != nil {
		return err
	}

	med.Holder = newHolder

	medAsBytes, _ := json.Marshal(med)

	return ctx.GetStub().PutState(medNumber, medAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create medstore chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting medstore chaincode: %s", err.Error())
	}
}
