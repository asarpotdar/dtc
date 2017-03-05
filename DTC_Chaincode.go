package main

import (
	"fmt"
	"errors"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Digital Trade Chain Chaincode
type DTCChaincode struct {
}

var contractIndexTxStr = "_contractIndexTxStr"
var buyerIndexTxStr = "_buyerIndexTxStr"

type ContractData struct{
	CONTRACT_ID string `json:"CONTRACT_ID"`
	CONTRACT_TITLE string `json:"CONTRACT_TITLE"`
	TOTAL_PRICE string `json:"TOTAL_PRICE"`
	CURRENCY string `json:"CURRENCY"`
	DELIVERY_DATE string `json:"DELIVERY_DATE"`
	ORDER_DATA string `json:"ORDER_DATA"`
	SELLER_ID string `json:"SELLER_ID"`
	BUYER_ID string `json:"BUYER_ID"`
	DEL_ADDR string `json:"BUYER_ID"`
	STATUS string `json:"STATUS"`

}

type Buyer struct{
	BUYER_ID string `json:"BUYER_ID"`
	BUYER_NAME string `json:"BUYER_NAME"`
	BUYER_BANK_ID string `json:"BUYER_BANK_ID"`
}

type Seller struct{
	SELLER_ID string `json:"SELLER_ID"`
	SELLER_NAME string `json:"SELLER_NAME"`
	SELLER_BANK_ID string `json:"SELLER_BANK_ID"`
}

// Init resets all the things
func (t *DTCChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 var err error
 if len(args) != 1 {
  return nil, errors.New("Incorrect number of arguments. Expecting 1")
 }
 if function == "InitContract"{
  t.InitContract(stub, args);
 }

 if err != nil {
  return nil, err
 }
 return nil, nil
}

func (t *DTCChaincode) InitContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 var err error

 fmt.Println("Init for Customer Data")
 var emptyContract []ContractData
 jsonAsBytes, _ := json.Marshal(emptyContract)
 err = stub.PutState(contractIndexTxStr, jsonAsBytes)
 if err != nil {
  return nil, err
 }
 return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *DTCChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 fmt.Println("invoke is running " + function)
 // register customer
 if function == "saveContract" {
  return t.saveContract(stub, args)
 }
 if function == "addBuyer" {
	return t.addBuyer(stub, args)
 }
 fmt.Println("invoke did not find func: " + function)     //error
 return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *DTCChaincode) saveContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var ContractDataObj ContractData
	var ContractDataList []ContractData
	var err error

	if len(args) != 7 {
		return nil, errors.New("Incorre		ct number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode
	ContractDataObj.CONTRACT_ID = args[0]
	ContractDataObj.CONTRACT_TITLE = args[1]
	ContractDataObj.TOTAL_PRICE = args[2]
	ContractDataObj.CURRENCY = args[3]
	ContractDataObj.SELLER_ID = args[4]
	ContractDataObj.BUYER_ID = args[5]
	ContractDataObj.DELIVERY_DATE = args[6]

	fmt.Printf("Input from user:%s\n", ContractDataObj)

	contractTxsAsBytes, err := stub.GetState(contractIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get contract Transactions")
	}
	json.Unmarshal(contractTxsAsBytes, &ContractDataList)

	ContractDataList = append(ContractDataList, ContractDataObj)
	jsonAsBytes, _ := json.Marshal(ContractDataList)

	err = stub.PutState(contractIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *DTCChaincode) addBuyer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var BuyerObj Buyer
	var err error

	if len(args) != 6 {
		return nil, errors.New("Incorre		ct number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode
	BuyerObj.BUYER_ID = args[0]
	BuyerObj.BUYER_NAME = args[1]
	BuyerObj.BUYER_BANK_ID = args[2]

	fmt.Printf("Input from user:%s\n", BuyerObj)

	jsonAsBytes, _ := json.Marshal(BuyerObj)

	err = stub.PutState(buyerIndexTxStr, jsonAsBytes)

	return jsonAsBytes, err
}


// Query callback representing the query of a chaincode
func (t *DTCChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

	var id string // Entities
	var err error
	var resAsBytes []byte

	id = args[0]

	if function == "GetContractDetails" {
		resAsBytes, err = t.GetContractDetails(stub, id)
  }
  if function == "getBuyers" {
		resAsBytes, err = t.GetBuyers(stub, id)
  }
	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	if err != nil {
	return nil, errors.New("Failed to Marshal the required Obj")
	}
	return resAsBytes, nil
}

func (t *DTCChaincode) GetContractDetails(stub shim.ChaincodeStubInterface, contractId string) ([]byte, error) {

	//var requiredObj RegionData
	var objFound bool
	ContractTxsAsBytes, err := stub.GetState(contractIndexTxStr)
	fmt.Printf("Output from chaincode >>>>> : %s\n", ContractTxsAsBytes)
	if err != nil {
	return nil, errors.New("Failed to Marshal the required Obj")
	}
	var ContractTxDataObjs []ContractData
	var ContractTxDataObjs1 []ContractData
	json.Unmarshal(ContractTxsAsBytes, &ContractTxDataObjs)
	length := len(ContractTxDataObjs)
	fmt.Printf("Output from chaincode: %s\n", ContractTxsAsBytes)

	if contractId == "" {
		res, err := json.Marshal(ContractTxDataObjs)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := ContractTxDataObjs[i]
		if contractId == obj.CONTRACT_ID {
			ContractTxDataObjs1 = append(ContractTxDataObjs1,obj)
			//requiredObj = obj
			objFound = true
		}
	}

	if objFound {
		res, err := json.Marshal(ContractTxDataObjs1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}

func (t *DTCChaincode) GetBuyers(stub shim.ChaincodeStubInterface, contractId string) ([]byte, error) {

	//var requiredObj RegionData
	BuyerTxsAsBytes, err := stub.GetState(buyerIndexTxStr)
	//if err != nil {
		//return nil, errors.New("Failed to get Merchant Transactions")
	//}
	//var BuyerObjs []Buyer
	res, err := json.Marshal("No Data found")
	fmt.Printf("Output from chaincode: %s\n", BuyerTxsAsBytes)
	return res, err

}


func main() {
	err := shim.Start(new(DTCChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
