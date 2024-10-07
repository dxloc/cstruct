package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func main() {
	test := []int{0, 0, 0, 1}

	type MyStruct struct {
		Value int32    `cstruct:"le"`
		Array [4]int16 `cstruct:"be"`
		Msg   string   `cstruct:"-"`
	}

	type MyStruct2 struct {
		Value2 int16    `cstruct:"le"`
		M      MyStruct `cstruct:"-"`
	}

	type MyStruct3 struct {
		Be int32       `cstruct:"le"`
		Le int32       `cstruct:"be"`
		A  []MyStruct2 `cstruct:"-"`
	}

	if test[0] != 0 {
		// Test 1: Struct with nested struct, last field is string
		fmt.Println("Test 1: Struct with nested struct")
		a := MyStruct2{
			Value2: 456,
			M: MyStruct{
				Value: 123,
				Array: [4]int16{1, 2, 3, 4},
				Msg:   "Hello, World!",
			},
		}
		fmt.Println(a) // print {456 {123 [1 2 3 4] Hello, World!}}

		b := cstruct.ToBytes(&a)
		fmt.Println(b) // print [200 1 123 0 0 0 0 1 0 2 0 3 0 4 72 101 108 108 111 44 32 87 111 114 108 100 33 0]

		var c MyStruct2
		cstruct.FromBytes(b, &c)
		fmt.Println(c) // print {456 {123 [1 2 3 4] Hello, World!}}
	}

	if test[1] != 0 {
		// Test 2: Struct with slice of struct, dynamic length struct, slice has 1 element
		fmt.Println("Test 2: Struct with slice of struct, dynamic length struct, slice has 1 element")
		a := MyStruct3{
			Be: 123, Le: 456, A: []MyStruct2{
				{Value2: 789, M: MyStruct{Value: 10, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 0"}},
			}}
		fmt.Println(a) // print {123 456 [{789 {10 [1 2 3 4] Hello, World! 0}}]}

		b := cstruct.ToBytes(&a)
		fmt.Println(b) // print [123 0 0 0 0 0 1 200 21 3 10 0 0 0 0 1 0 2 0 3 0 4 72 101 108 108 111 44 32 87 111 114 108 100 33 32 48 0]

		var c MyStruct3
		cstruct.FromBytes(b, &c)
		fmt.Println(c) // print {123 456 [{789 {10 [1 2 3 4] Hello, World! 0}}]}
	}

	if test[2] != 0 {
		// Test 3: Struct with slice of struct, contains more than 1 element
		// This will fail to deserialize because len(aa.A) > 1 and MyStruct2 is an unknown length struct.
		fmt.Println("Test 3: Struct with slice of struct with more than 1 element")
		a := MyStruct3{
			Be: 123, Le: 456, A: []MyStruct2{
				{Value2: 789, M: MyStruct{Value: 10, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 0"}},
				{Value2: 123, M: MyStruct{Value: 11, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 1"}},
			}}
		fmt.Println(a) // print {123 456 [{789 {10 [1 2 3 4] Hello, World! 0}} {123 {11 [1 2 3 4] Hello, World! 1}}]}

		b := cstruct.ToBytes(&a)
		fmt.Println(b) // print [123 0 0 0 0 0 1 200 21 3 10 0 0 0 0 1 0 2 0 3 0 4 123 0 11 0 0 0 0 1 0 2 0 3 0 4]

		var c MyStruct3
		cstruct.FromBytes(b, &c)
		fmt.Println(c) // {123 456 [{789 {10 [1 2 3 4] }} {123 {11 [1 2 3 4] }}]}
	}

	if test[3] != 0 {
		// Test 4: Struct with nested struct, fixed length
		fmt.Println("Test 4: Struct with nested struct, fixed length")
		type MyStruct4 struct {
			Value1 int32 `cstruct:"le"`
			Value2 int32 `cstruct:"le"`
		}

		type MyStruct5 struct {
			Value3 int32     `cstruct:"le"`
			M      MyStruct4 `cstruct:"-"`
			Value4 int32     `cstruct:"le"`
		}

		a := MyStruct5{
			Value3: 456,
			M: MyStruct4{
				Value1: 123,
				Value2: 456,
			},
			Value4: 789,
		}
		fmt.Println(a) // print {456 {123 456} 789}
		b := cstruct.ToBytes(&a)
		fmt.Println(b) // print

		var c MyStruct5
		cstruct.FromBytes(b, &c)
		fmt.Println(c) // print {456 {123 456} 789}
	}
}
