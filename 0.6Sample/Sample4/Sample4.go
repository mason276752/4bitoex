package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Sample4 struct {
}

func (this *Sample4) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	stub.CreateTable("SampleTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "colBool", Type: shim.ColumnDefinition_BOOL, Key: true},
		&shim.ColumnDefinition{Name: "colBytes", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "colString", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "colInt32", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "colUint64", Type: shim.ColumnDefinition_UINT64, Key: true},
	})
	//新增五筆相似的資料用於示範
	stub.InsertRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: true}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Bytes 1")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 1}},
		},
	})
	stub.InsertRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: true}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Bytes 2")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 2}},
		},
	})
	stub.InsertRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: true}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Bytes 3")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 3}},
		},
	})
	stub.InsertRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: false}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Bytes 4")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 2}},
		},
	})
	//重複key值 新增失敗
	stub.InsertRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: true}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Bytes 5")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 1}},
		},
	})
	return nil, nil
}

func (this *Sample4) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//欄位更新:比對相同key值,更新非key欄位
	//故Key為true之欄位不可更改
	ok, err := stub.ReplaceRow("SampleTable", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_Bool{Bool: false}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte("This Is Sample")}},
			&shim.Column{Value: &shim.Column_String_{String_: "This Is String"}},
			&shim.Column{Value: &shim.Column_Int32{Int32: rand.Int31()}},
			&shim.Column{Value: &shim.Column_Uint64{Uint64: 2}},
		},
	})
	if err != nil {
		shim.NewLogger("My Error").Infof("Error:%v\n", err.Error())
		return nil, err
	}
	if ok {
		fmt.Printf("OK")
	}

	return nil, nil
}

func (this *Sample4) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var cols []shim.Column
	if (len(args)) > 0 {
		key1, _ := strconv.ParseBool(args[0])
		cols = append(cols, shim.Column{Value: &shim.Column_Bool{Bool: key1}})
	}
	if (len(args)) > 1 {
		key2 := args[1]
		cols = append(cols, shim.Column{Value: &shim.Column_String_{String_: key2}})
	}
	var key3 uint64
	if (len(args)) > 2 {
		key3, _ = strconv.ParseUint(args[2], 10, 0)
		cols = append(cols, shim.Column{Value: &shim.Column_Uint64{Uint64: key3}})
	}

	switch function {
	case "row":
		//只能比對相等,不能比對不相等大於或小於
		//必須每個key都必填 才能順利搜尋目標
		//只會找到一個
		row, err := stub.GetRow("SampleTable", cols)
		if err != nil {
			shim.NewLogger("My Error").Infof("Error:%v\n", err.Error())
			return nil, err
		}
		shim.NewLogger("Info").Infof("%v\n", row)
		break
	//----------------------------------------------
	case "rows":
		//不可填滿每個key,否則搜尋無結果
		//只能依序比對key值,不可跳過前面的key欄位
		rows, err := stub.GetRows("SampleTable", cols)
		if err != nil {
			fmt.Printf(err.Error())
			shim.NewLogger("My Error").Infof("Error:%v\n", err.Error())
			return nil, err
		}
		for i := 1; ; i++ {
			row, ok := <-rows
			if !ok {
				rows = nil
			} else {
				fmt.Printf("%d:%v\n", i, row)
			}
			if rows == nil {
				break
			}
		}
		break
	}
	return []byte("Query Succeed"), nil
}

func main() {
	err := shim.Start(new(Sample4))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
	fmt.Print("Sample4")
}
