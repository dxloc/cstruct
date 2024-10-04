package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func main() {
	type MyStruct struct {
		Value int32    `cstruct:"le"`
		Array [4]int32 `cstruct:"be"`
		Msg   string   `cstruct:"-"`
	}

	type MyStruct2 struct {
		Value int32    `cstruct:"le"`
		M     MyStruct `cstruct:"-"`
	}

	type MyStruct3 struct {
		Be int32   `cstruct:"le"`
		Le int32   `cstruct:"be"`
		A  []int32 `cstruct:"be"`
	}

	a := MyStruct2{
		Value: 456,
		M: MyStruct{
			Value: 123,
			Array: [4]int32{1, 2, 3, 4},
			Msg:   "Hello, World!",
		},
	}
	fmt.Println(a) // print {456 {123 [1 2 3 4] Hello, World!}}

	b := cstruct.ToBytes(&a)
	fmt.Println(b) // print [200 1 0 0 123 0 0 0 0 0 0 1 0 0 0 2 0 0 0 3 0 0 0 4 72 101 108 108 111 44 32 87 111 114 108 100 33 0]

	var c MyStruct2
	cstruct.FromBytes(b, &c)

	fmt.Println(c) // print {456 {123 [1 2 3 4] Hello, World!}}

	aa := MyStruct3{Be: 123, Le: 456, A: []int32{789, 10}}
	fmt.Println(aa) // print {123 456 [789 10]}

	bb := cstruct.ToBytes(&aa)
	fmt.Println(bb) // print [123 0 0 0 0 0 1 200 0 0 3 21 0 0 0 10]

	var cc MyStruct3
	cstruct.FromBytes(bb, &cc)
	fmt.Println(cc) // print {123 456 [789 10]}
}
