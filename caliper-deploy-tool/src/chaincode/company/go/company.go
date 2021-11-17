/*
提交表单至区块链
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//err
var err error

type Company struct {
	_id 				int 			`json:"_id"`
	History 			[]History 		`json:"history"` 
	Href 				string 			`json:"href"` 
	Short_name 			string 			`json:"short_name"` 
	Company_word		string   		`json:"company_word"` 
	Basic				Basic			`json:"basic"`
	Identification		int				`json:"identification"`
	Logo				string			`json:"logo"`
	Company_img			[]string		`json:"company_img"`
	Name				string			`json:"name"`
	Lagou_url			string			`json:"lagou_url"`
	Manager_list 		[]Manager_list	`json:"manager_list"`
	Address				[]string		`json:"address"`
	Company_intro_text 	string			`json:"company_intro_text"`
}

type History struct {
	Url  	string 		`json:"url"`
	Date  	string      `json:"date"`
	Day 	string		`json:"day"`
	Title  	string		`json:"title"`
	Content string 		`json:"content"`
	Type 	string		`json:"type"`
}

type Basic struct {
	Process  string 	`json:"process"`
	Type  	 string     `json:"type"`
	Address  string		`json:"address"`
	Number   string		`json:"number"`
}

type Manager_list struct{
	Content  string 	`json:"content"`
	Weibo  	 string     `json:"weibo"`
	Title    string		`json:"title"`
	Name     string		`json:"name"`
}



//SimpleChaincode Init方法
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("----- company Init")
	return shim.Success(nil)
}

//invoke,func,
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("----- company Invoke")
	function, args := stub.GetFunctionAndParameters()

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting at least 1")
	}
	//不同的函数调用
	switch {
	case function == "submit_company": //提交的接收调拨单，将调拨单记入区块链
		return t.submit_company(stub, args)
	case function == "query":
		return t.query(stub, args)
	default:
		fmt.Printf("function is not exist\n")
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}


func (t *SimpleChaincode) submit_company(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("submit_company")
	var companyInfo Company
	fmt.Println(args)
	fmt.Println(args[0])
	fmt.Println(args[1])


	// str_params, err := base64.StdEncoding.DecodeString(args[0])
	err = json.Unmarshal([]byte(args[1]), &companyInfo)


	// fmt.Printf("lmxlmxlmxlmlmx %s",companyInfo)
	companyShortName := args[0]

	// companyBytes, err := stub.GetState(companyShortName)
	// if err != nil {
	// 	return shim.Error("Failed to get company:" + err.Error())
	// } else if companyBytes != nil {
	// 	return shim.Error("This bill already exists:" + companyShortName)
	// }

	json_company, err := json.Marshal(companyInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(companyShortName, json_company)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("submit_company success \n", companyShortName)
	return shim.Success(json_company)
}


func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("----- company query  \n")

	AllocationBillID := args[0]

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	AllocationB, err := stub.GetState(AllocationBillID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get AllocationBill for\"}" + AllocationBillID + "\"}"
		return shim.Error(jsonResp)
	}
	if AllocationB == nil {
		jsonResp := "{\"Error\":\"Nil amount for\"}" + AllocationBillID + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"AllocationBillID\":\"" + AllocationBillID + "\",\"Bill\":\"" + string(AllocationB) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(AllocationB)

}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
