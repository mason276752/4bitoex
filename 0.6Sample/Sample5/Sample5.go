package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

type Sample5 struct {
}

var ccName = "Sample2"

func (this *Sample5) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//創建表格用來記錄智能合約事件
	stub.CreateTable("ccTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "User", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "point", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "x", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "y", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "m", Type: shim.ColumnDefinition_INT32, Key: false},
		&shim.ColumnDefinition{Name: "unixtime", Type: shim.ColumnDefinition_INT64, Key: false},
	})

	return nil, nil
}

func (this *Sample5) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var x, y, m, point int32
	var unixtime int64

	switch function {
	case "getPoint":
		var cols []shim.Column
		cols = append(cols, shim.Column{Value: &shim.Column_String_{String_: args[0]}})
		//查詢智能合約事件表
		row, err := stub.GetRow("ccTable", cols)
		if err != nil {
			shim.NewLogger("Error").Errorf("%v", err.Error())
			return nil, err
		}
		//比對是否時間已到
		if time.Now().Unix() > row.Columns[5].GetInt64() {
			if row.Columns[2].GetInt32() != 0 {
				if len(args) == 3 {
					x, err := strconv.Atoi(args[1])
					if err != nil {
						shim.NewLogger("Error").Errorf("%v", err.Error())
						return nil, err
					}

					y, err := strconv.Atoi(args[2])
					if err != nil {
						shim.NewLogger("Error").Errorf("%v", err.Error())
						return nil, err
					}
					//比對地點是否已達
					m := math.Pow(float64(row.Columns[2].GetInt32()-int32(x)), 2) + math.Pow(float64(row.Columns[3].GetInt32()-int32(y)), 2)
					if float64(row.Columns[4].GetInt32()) >= math.Sqrt(m) {
						//修改受予人金額
						funcAndArgs := util.ToChaincodeArgs("get", args[0])
						response, _ := stub.QueryChaincode(ccName, funcAndArgs)
						temp, _ := strconv.Atoi(string(response))
						temp += int(row.Columns[1].GetInt32())
						funcAndArgs = util.ToChaincodeArgs("put", args[0], strconv.Itoa(int(temp)))
						stub.InvokeChaincode(ccName, funcAndArgs)
						//刪除此智能合約事件
						stub.DeleteRow("ccTable", cols)
					}
				}
			}
		}

		break
	case "sendPoint":
		m = 10
		i := 0
		for ; i < len(args); i++ {
			if _, err := strconv.Atoi(args[i]); err == nil {
				temp, _ := strconv.Atoi(args[i])
				x = int32(temp)
				temp, _ = strconv.Atoi(args[i+1])
				y = int32(temp)
				i += 1
			} else if args[i] == "m" {
				temp, _ := strconv.Atoi(args[i+1])
				m = int32(temp)
				i += 1
			} else if args[i] == "time" {
				temp, _ := strconv.Atoi(args[i+1])
				unixtime = time.Now().Unix() + int64(temp*60)
				i += 1
			} else if args[i] == "point" {
				temp, _ := strconv.Atoi(args[i+1])
				point = int32(temp)
				i += 1

			} else if args[i] == "from" {
				//修改給予人的金額
				fromuser := args[i+1]
				funcAndArgs := util.ToChaincodeArgs("get", fromuser)
				response, _ := stub.QueryChaincode(ccName, funcAndArgs)
				temp, _ := strconv.Atoi(string(response))
				temp -= int(point)
				funcAndArgs = util.ToChaincodeArgs("put", fromuser, strconv.Itoa(int(temp)))
				info, err := stub.InvokeChaincode(ccName, funcAndArgs)
				if err != nil {
					shim.NewLogger("Error").Errorf("%v", err.Error())
					return nil, err
				}
				shim.NewLogger("Info").Infof("%v", info)
				i += 1
			}

		}
		i = i - 1
		//新增一筆智能合約事件
		stub.InsertRow("ccTable", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: args[i]}},
				&shim.Column{Value: &shim.Column_Int32{Int32: point}},
				&shim.Column{Value: &shim.Column_Int32{Int32: x}},
				&shim.Column{Value: &shim.Column_Int32{Int32: y}},
				&shim.Column{Value: &shim.Column_Int32{Int32: m}},
				&shim.Column{Value: &shim.Column_Int64{Int64: unixtime}},
			},
		})
		break
	}
	return nil, nil
}

//純查詢Sample2
func (this *Sample5) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	funcAndArgs := util.ToChaincodeArgs("get", args[0])
	response, _ := stub.QueryChaincode(ccName, funcAndArgs)
	return response, nil
}

func main() {
	err := shim.Start(new(Sample5))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
	fmt.Print("Sample5")
}
