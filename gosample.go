package main

import (
	"fmt"
)

var i = 5

func add(val *int) (int, int) { //指標參數和多回傳值
	*val++
	return *val % 2, *val * *val
}
func main() {
	var str []string
	str = append(str, "Hello", "World", "Hello", "Sample")
	for index, value := range str[1:] { //範圍的索引&遍歷陣列
		fmt.Println("index:", index, "value:", value)
	}
	for i > 0 { //while(i>0)
		i--
		fmt.Println("i:", i)
	}
	for { //無限迴圈
		if mod, square := add(&i); mod == 0 { //多回傳和if合併
			if i > 0 {
				fmt.Println("i:", i, "square:", square)
			}
			if square > 1000 {
				break
			}
		}
	}

	sp := new(sample)
	sp.setN(i)
	fmt.Println("sp's n:", sp.getN())
}

type sample struct {
	n int
}

func (this *sample) setN(_n int) {
	this.n = _n

}
func (this *sample) getN() int {
	return this.n
}
