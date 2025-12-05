package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func init() {}

func main() {
	type MyStruct struct {
		Arr []int32 `cstruct:"be"`
	}

	a := MyStruct{
		Arr: []int32{1, 2, 3},
	}
	fmt.Println(a)

	b := cstruct.Marshal(&a)
	fmt.Println(b)

	var c MyStruct
	cstruct.Unmarshal(b, &c)
	fmt.Println(c)
}
