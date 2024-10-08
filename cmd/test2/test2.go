package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func init() {}

func main() {
	type MyStruct struct {
		Id      uint32 `cstruct:"be"`
		Action  int32  `cstruct:"be"`
		Content []byte `cstruct:"-"`
	}

	a := MyStruct{
		Id:      123,
		Action:  456,
		Content: []byte("Hello, World!"),
	}
	fmt.Println(a)

	b := cstruct.ToBytes(&a)
	fmt.Println(b)

	var c MyStruct
	cstruct.FromBytes(b, &c)
	fmt.Println(c)
}
