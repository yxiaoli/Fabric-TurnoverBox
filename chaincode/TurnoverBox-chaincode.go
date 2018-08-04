// SPDX-License-Identifier: Apache-2.0

/*
 Author: Sherry
 Time: 2018 
 This code is based on code written by the Hyperledger Fabric community.
*/

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
  formatting, and string manipulation  
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts  
*/ 
import (
	"bytes"
	"encoding/json" //
	"fmt"
	"strconv"
    "container/list" 
    "time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Box structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Box struct {
	BOwner string `json:"name"`
	Start string `json:"S_timestamp"`
	Type string `json:"Type"`
	End string `json:"E_timestamp"`
}

/* Define Cash/Coin structure, with only one property.  
Structure tags are used by encoding/json library
*/
type Coin struct{
	COwner string `json:"name"`
}

/*
 * The Init method *
 called when the Smart Contract "TurnoverCh-chaincode" is instantiated by the *network*
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "TurnoverBox-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryBox" {
		return s.queryBox(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "addBox" {
		return s.addBox(APIstub, args)
	} else if function == "queryAllBox" {
		return s.queryAllBox(APIstub)
	} else if function == "refuelFee" {
		return s.refuelFee(APIstub, args)
	} else if function == "depositCoin" {
		return s.depositCoin(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryBox method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
 */
func (s *SmartContract) queryBox(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	return shim.Success(tunaAsBytes)
}

/*
 * The initLedger method *
Will add test data to our network
this method is only temporarily
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	boxes := []Box{
		Box{BOwner: "Operator", Start: "20180314001", Type: "Box-N", End "INF"},
		Box{BOwner: "Operator", Start: "20180314002", Type: "Box-N", End "INF"},
		Box{BOwner: "Supplier", Start: "20180314001", Type: "Box-A", End "20180621002"},
		Box{BOwner: "Distributor", Start: "20180714001", Type: "Box-N", End "20180921002"}		
	}
    coins := []COIN{
    	Coin{COwner: "Supplier"},
    	Coin{COwner: "Supplier"},
  	 	Coin{COwner: "Distributor"}
    	Coin{COwner: "Retailer"}
   
    }
	i := 0
	for i < len(boxes) {
		fmt.Println("i is ", i)
		tunaAsBytes, _ := json.Marshal(boxes[i])
		coinAsBytes, _ := json.Marshal(coins[i])		
		APIstub.PutState(strconv.Itoa(i+1), tunaAsBytes)
		APIstub.PutState(strconv.Itoa(i+1), coinAsBytes)
		fmt.Println("Added", tuna[i], coins[i])
		i = i + 1
	}

     
	return shim.Success(nil)
}

/*
	 * The addBox method *
The Box operator will add new turnover boxes into the network

 */
func (s *SmartContract) addBox(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
    if args[1] != "Operator"{
    	return shim.Error("The Only Operator can add Boxes")

    }
	var box = Box{ Owner: args[1], Start: args[2], Type: args[3], End: args[4] }

	boxAsBytes, _ := json.Marshal(box)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add box record : %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllBox method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllBox(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTuna:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The refuelFee method *
 The supplier should pay for the
//This function takes in 2 arguments, number of boxes and Boxes Type. 
 */
func (s *SmartContract) refuelFee(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

    Bnum := args[0]
    Btype := args[1]
  	queryString_b := fmt.Sprintf("{\"selector\":{\"BOwner\":\"Operator\",\"Type\":\"%s\"}}", Btype)
  	queryString_c := fmt.Sprintf("{\"selector\":{\"COwner\":\"Supplier\"}}")
  	
  	resultsIterator_b, err_b := APIstub.GetQueryResult(queryString_b)
  	resultsIterator_c, err := APIstub.GetQueryResult(queryString_c)
    if err_b != nil {
		return shim.Error(err_b.Error())
	}

 	if err_c != nil {
		return shim.Error(err_c.Error())
	}

	boxes := list.New()
	coins := list.New()
    i := 0
    t := time.Now()
	for i < Bnum {
     
        if(!resultsIterator_b.HasNext()){
        	return shim.Error("Insufficient Box, only %d boxe(s) has been pledged",i)
        }

        if(!resultsIterator_c.HasNext()){
        	return shim.Error("Insufficient Money, only %d boxe(s) has been pledged",i)
        }
	    queryResponse_b, err_b := resultsIterator_b.Next()
		queryResponse_c, err_c := resultsIterator_c.Next()
		
		if err_b != nil {
			return nil, err_b
		}
		if err_c != nil {
			return nil, err_c
		}
        // change the owner of box
        key_b = queryResponse_b.Key
        value_b = queryResponse_b.Value

		value_b.Start = t.String()
		value_b.End = t.AddDate(0,1,0)
		value_b.Owner = "Supplier"
       	boxAsBytes, _ = json.Marshal(value_b)

        PutState(key_b, boxAsBytes)
        //change the owner of coins
        key_c = queryResponse_c.Key
        value_c = queryResponse_c.Value
		value_c.Owner = "Operator"
       	coinAsBytes, _ = json.Marshal(value_c)

        PutState(key_b, coinAsBytes)        
		i ++
	}
	
	return shim.Success(nil)
}

/*
 *  The depositCoin method *
 the clients except the operator should deposit some coins
 in order to pay pledge 
 the func takes in 2 arguments, owner, num of coins 
*/
 
func (s *SmartContract) depositCoin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response{
 	
 	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

    var num = args[1]
 
    for i := 0; i< num; i++ {
	   var coins = Coin{ Owner: args[0] }
	   coinAsBytes,_ := json.Marshal(coins)	
	   err := APIstub.PutState(args[0], coinAsBytes)
	   if err != nil {
	   		return shim.Error(fmt.Sprintf("Failed to add box record : %s", args[0]))
	   }
	}
	   
	return shim.Success(nil)
}



/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}