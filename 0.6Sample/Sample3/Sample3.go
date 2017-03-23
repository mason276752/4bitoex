package main

import (
	"fmt"

	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

type Sample3 struct {
}

func (this *Sample3) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

var ccName = "Sample2"

//util通常只會用到ToChaincodeArgs和ArrayToChaincodeArgs
func (this *Sample3) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//這邊試著呼叫 Sample2 的chaincode
	//調用 invoke
	funcAndArgs := util.ToChaincodeArgs("put", "Alice", "300")
	response, err := stub.InvokeChaincode(ccName, funcAndArgs)
	if err != nil {
		shim.NewLogger("Error").Errorf("%v", err.Error())
		return nil, err
	}
	shim.NewLogger("Invoke Response").Infof("%v\n", string(response))

	//調用 query
	funcAndArgs = util.ToChaincodeArgs("get", "Alice")
	response, err = stub.QueryChaincode(ccName, funcAndArgs)
	if err != nil {
		shim.NewLogger("Error").Errorf("%v", err.Error())
		return nil, err
	}
	shim.NewLogger("Query Response").Infof("%v\n", string(response))
	return nil, nil
}

func (this *Sample3) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//調用 invoke

	funcAndArgs := util.ToChaincodeArgs("put", "Alice", "300")
	response, err := stub.InvokeChaincode(ccName, funcAndArgs)
	if err != nil {
		//這會發生錯誤query中也不能調用invoke
		shim.NewLogger("Error").Errorf("query中不能調用invoke:%v", err.Error())
	}
	//調用1 query1
	funcAndArgs = util.ToChaincodeArgs("notput", "Alice", "450")
	response, err = stub.QueryChaincode(ccName, funcAndArgs)
	if err != nil {
		//這會返回Sample2的錯誤
		shim.NewLogger("Error").Errorf("%v", err.Error())
	}
	//調用 query2
	funcAndArgs = util.ToChaincodeArgs("get", "Alice")
	response, err = stub.QueryChaincode(ccName, funcAndArgs)
	if err != nil {
		shim.NewLogger("Error").Errorf("%v", err.Error())
		return nil, errors.New("Error")
	}
	shim.NewLogger("Query Response").Infof("Query:%v", string(response))
	return response, nil
}

func main() {
	err := shim.Start(new(Sample3))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
	fmt.Print("Sample3")
}
