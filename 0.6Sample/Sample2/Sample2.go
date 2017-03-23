package main

import (
	"fmt"
	"math/rand"

	"strconv"

	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Sample2 struct {
}

//當寫入,刪除資料世界狀態也會隨之改變
func (this *Sample2) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//stub 是作為讀寫區塊鏈,調用chaincode,查詢交易相關消息,簽章驗證有關的物件
	//建立或覆寫一個key-value
	err := stub.PutState("Alice", []byte(args[0]))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("Bob", []byte("10000"))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//當寫入,刪除資料世界狀態也會隨之改變
func (this *Sample2) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "put" {
		//建立或復寫一個key-value
		stub.PutState(args[0], []byte(args[1]))
	} else if function == "del" {
		//刪除一個key-value
		stub.DelState(args[0])
	}
	//返回隨機值Sample3用到
	return []byte("Sample2 Invoke,Rand:" + strconv.Itoa(rand.Int())), nil
}

//當寫入,刪除資料會發生錯誤,所以不能這樣做
//故世界狀態不會被改變
func (this *Sample2) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	if function == "get" {
		//讀取一個key-value
		value, err := stub.GetState(args[0])
		if err != nil {
			return nil, err
		}
		return value, nil
	} else {
		//建立或復寫一個key-value
		//必定失敗
		err = stub.PutState(args[0], []byte(args[1]))
		if err != nil {
			log := shim.NewLogger("This Is Log")
			log.Error("Query中不可PutState")
			return nil, errors.New("[Sample2]錯誤返回")

		}
		return nil, errors.New("[Sample2]正常返回")
	}

}

func main() {
	err := shim.Start(new(Sample2))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
	fmt.Print("Sample2")
}
