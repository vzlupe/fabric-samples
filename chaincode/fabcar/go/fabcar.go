/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 5 utility libraries for formatting, handling bytes, reading and writing JSON, string manipulation, and math operations
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"math"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the temp sensor structure
type TempSensor struct {
        Temp string `json:"temp"`
        Peer0 string `json:"peer0"`
        Peer1 string `json:"peer1"`
        Peer2 string `json:"peer2"`
        Resp string `json:"resp"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTemp" {
		return s.queryTemp(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createTemp" {
		return s.createTemp(APIstub, args)
	} else if function == "queryAllTemps" {
		return s.queryAllTemps(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryTemp(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tempAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(tempAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	temps := []TempSensor{
                TempSensor{Temp: "0",Peer0: "none",Peer1: "none",Peer2: "none",Resp: "none"},
	}

	//initialize the ledger with a blank record
	i := 0
	for i < len(temps) {
		fmt.Println("i is ", i)
		tempAsBytes, _ := json.Marshal(temps[i])
		APIstub.PutState("TEMP"+strconv.Itoa(i), tempAsBytes)
		fmt.Println("Added", temps[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createTemp(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	//declare variables for the peers individual decision and the group response, as well as the acceptable temp delta
	var res = "none"
	var dec = "none"
	var acp = 1.1

	//get the number of the transaction e.g. 'TEMP5' -> 5
	num, _ := strconv.Atoi(args[0][4:])

	//if this is the first transaction, accept. if not, query the blockchain for the previous transaction, compare
	//its temp with the new temp, and accept if the delta is acceptable reject if not
	if (num > 1){
		var checkKey = "TEMP" + strconv.Itoa(num - 1)
		t := TempSensor{}
		tempAsBytes2, _ := APIstub.GetState(checkKey)
		json.Unmarshal(tempAsBytes2, &t)
		otemp, _ := strconv.ParseFloat(t.Temp, 64)
		ntemp, _ := strconv.ParseFloat(args[1], 64)
		if (math.Abs(otemp - ntemp) > acp){
			dec = "reject"
		} else {
			dec = "accept"
		}
	} else {
		dec = "accept"
	}

	//compare the peer's decision against those of the 'false peers'
	//if they all agree, accept. otherwise, reject with dissenting peer
	if (dec != args[2] || dec != args[3] || args[2] != args[3]) {
		if (dec == args[2]){
			res = "unsuccessful-peer2"
		} else if (dec == args[3]) {
			res = "unsuccessful-peer1"
		} else {
			res = "unsuccessful-peer0"
		}
	} else {
		res = "successful"
	}


	//set the full transaction request and append to the blockchain
	var temp = TempSensor{Temp: args[1],Peer0: dec,Peer1: args[2],Peer2: args[3],Resp: res}
	tempAsBytes, _ := json.Marshal(temp)
	APIstub.PutState(args[0], tempAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllTemps(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TEMP0"
	endKey := "TEMP999"

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
		// Add a comma before array members, suppress it for the first array member
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

	fmt.Printf("- queryAllTemps:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
