package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type BenchmarkChaincode struct {
}

func (t *BenchmarkChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments. Expecting 3, got %v", len(args)))
	}

	// Initialize the chaincode
	var err error
	firstAccount, err := strconv.Atoi(args[0])
	accountCount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for number of accounts")
	}
	initValue, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for default account value")
	}
	// Write the state to the ledger

	for i := firstAccount; i < firstAccount+accountCount; i++ {
		err = stub.PutState("account"+strconv.Itoa(i), []byte(strconv.Itoa(initValue)))
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}

func (t *BenchmarkChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "transfer" {
		// Make payment of X units from A to B
		return t.transfer(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	} else if function == "init" {
		// the old "Query" is now implemtned in invoke
		return t.Init(stub)
	}

	return shim.Error(fmt.Sprintf("Invalid invoke function name. Expecting \"init\", \"transfer\" or \"query\", got \"%s\"", function))
}

// Transaction makes payment of X units from A to B
func (t *BenchmarkChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments. Expecting 3, got %v", len(args)))
	}

	A = args[0]
	B = args[1]

	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *BenchmarkChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(BenchmarkChaincode))
	if err != nil {
		fmt.Printf("Error starting benchmark chaincode: %s", err)
	}
}
