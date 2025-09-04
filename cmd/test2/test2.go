package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func init() {}

func main() {
	type MyStruct struct {
		A [4]byte `cstruct:"le"`
		B [4]byte `cstruct:"le"`
		C uint32  `cstruct:"le"`
	}

	a := MyStruct{
		A: [4]byte{1, 2, 3, 4},
		B: [4]byte{5, 6, 7, 8},
		C: 0x12345678,
	}
	fmt.Println(a)

	b := cstruct.Marshal(&a)
	fmt.Println(b)

	var c MyStruct
	cstruct.Unmarshal(b, &c)
	fmt.Println(c)
}
