package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim" //inculde 這檔案
)

//golang的結構相當於 c,java 結構,物件 [必要]
type Sample1 struct {
}

//deploy 所執行的函數[必要]
func (this *Sample1) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("function:%s\nargs:%#v", function, args)
	return []byte("[Sample1]deploy成功之訊息"), nil
}

//invoke 所執行的函數[必要]
func (this *Sample1) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("function:%s\nargs:%#v", function, args)
	return nil, errors.New("[Sample1]invoke失敗之訊息")
}

//query 所執行的函數[必要]
func (this *Sample1) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("function:%s\nargs:%#v", function, args)
	return []byte("[Sample1]查詢完成之返回訊息"), nil
}

//程式進入點
func main() {
	//main函數暫停在這行,等待與peer互動,直到shim意外中斷
	err := shim.Start(new(Sample1))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

	fmt.Print("Sample1")
}
