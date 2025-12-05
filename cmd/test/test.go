package main

import (
	"fmt"

	"github.com/dxloc/cstruct"
)

func main() {
	test := []int{0, 0, 0, 0, 0, 1}

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

	type MyStruct4 struct {
		Value1 int32
		Value2 int32
	}

	type MyStruct5 struct {
		Value3 int32     `cstruct:"le"`
		M      MyStruct4 `cstruct:"be"`
		Value4 int32     `cstruct:"le"`
	}

	type MyStruct6 struct {
		Be int32 `cstruct:"be"`
		Le int16 `cstruct:"le"`
	}

	type MyStruct7 struct {
		Value int32        `cstruct:"le"`
		A     [4]MyStruct6 `cstruct:"-"`
		B     string
	}

	type ArrMyStruct4 [2]MyStruct4

	type MyStruct8 struct {
		Nt0 int32
		Nt1 int32
		A   []ArrMyStruct4 `cstruct:"be"`
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
		fmt.Println(a) // prints {456 {123 [1 2 3 4] Hello, World!}}

		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [200 1 123 0 0 0 0 1 0 2 0 3 0 4 72 101 108 108 111 44 32 87 111 114 108 100 33 0]

		var c MyStruct2
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // prints {456 {123 [1 2 3 4] Hello, World!}}
	}

	if test[1] != 0 {
		// Test 2: Struct with slice of struct, dynamic length struct, slice has 1 element
		// The string won't be serialized
		fmt.Println("Test 2: Struct with slice of struct, dynamic length struct, slice has 1 element")
		a := MyStruct3{
			Be: 123, Le: 456, A: []MyStruct2{
				{Value2: 789, M: MyStruct{Value: 10, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 0"}},
			}}
		fmt.Println(a) // prints {123 456 [{789 {10 [1 2 3 4] Hello, World! 0}}]}

		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [123 0 0 0 0 0 1 200 21 3 10 0 0 0 0 1 0 2 0 3 0 4]

		var c MyStruct3
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // prints {123 456 [{789 {10 [1 2 3 4] }}]}
	}

	if test[2] != 0 {
		// Test 3: Struct with slice of struct, contains more than 1 element
		// The string won't be serialized
		fmt.Println("Test 3: Struct with slice of struct with more than 1 element")
		a := MyStruct3{
			Be: 123, Le: 456, A: []MyStruct2{
				{Value2: 789, M: MyStruct{Value: 10, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 0"}},
				{Value2: 123, M: MyStruct{Value: 11, Array: [4]int16{1, 2, 3, 4}, Msg: "Hello, World! 1"}},
			}}
		fmt.Println(a) // prints {123 456 [{789 {10 [1 2 3 4] Hello, World! 0}} {123 {11 [1 2 3 4] Hello, World! 1}}]}

		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [123 0 0 0 0 0 1 200 21 3 10 0 0 0 0 1 0 2 0 3 0 4 123 0 11 0 0 0 0 1 0 2 0 3 0 4]

		var c MyStruct3
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // {123 456 [{789 {10 [1 2 3 4] }} {123 {11 [1 2 3 4] }}]}
	}

	if test[3] != 0 {
		// Test 4: Struct with nested struct, fixed length
		fmt.Println("Test 4: Struct with nested struct, fixed length")

		a := MyStruct5{
			Value3: 456,
			M: MyStruct4{
				Value1: 123,
				Value2: 456,
			},
			Value4: 789,
		}
		fmt.Println(a) // prints {456 {123 456} 789}
		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [200 1 0 0 123 0 0 0 200 1 0 0 21 3 0 0]

		var c MyStruct5
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // prints {456 {123 456} 789}
	}

	if test[4] != 0 {
		// Test 5: Struct with array of struct
		fmt.Println("Test 5: Struct with array of string")

		a := MyStruct7{
			Value: 123,
			A: [4]MyStruct6{
				{Be: 789, Le: 10},
				{Be: 123, Le: 11},
				{Be: 456, Le: 12},
				{Be: 789, Le: 13},
			},
			B: "Hello, World!",
		}
		fmt.Println(a) // prints {123 [{789 10} {123 11} {456 12} {789 13}] Hello, World!}
		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [123 0 0 0 0 0 3 21 10 0 0 0 0 123 11 0 0 0 1 200 12 0 0 0 3 21 13 0 72 101 108 108 111 44 32 87 111 114 108 100 33 0]

		var c MyStruct7
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // prints {123 [{789 10} {123 11} {456 12} {789 13}] Hello, World!}
	}

	if test[5] != 0 {
		// Test 6: Slice of array of struct
		fmt.Println("Test 6: Slice of array of struct")

		a := MyStruct8{
			Nt0: 15,
			Nt1: 33,
			A: []ArrMyStruct4{
				[2]MyStruct4{
					{
						Value1: 789,
						Value2: 10,
					},
					{
						Value1: 123,
						Value2: 11,
					},
				},
				[2]MyStruct4{
					{
						Value1: 456,
						Value2: 12,
					},
					{
						Value1: 789,
						Value2: 13,
					},
				},
			},
		}
		fmt.Println(a) // prints {15 33 [[{789 10} {123 11}] [{456 12} {789 13}]]}
		b := cstruct.Marshal(&a)
		fmt.Println(b) // prints [15 0 0 0 33 0 0 0 21 3 0 0 10 0 0 0 123 0 0 0 11 0 0 0 200 1 0 0 12 0 0 0 21 3 0 0 13 0 0 0]

		var c MyStruct8
		cstruct.Unmarshal(b, &c)
		fmt.Println(c) // prints {15 33 [[{789 10} {123 11}] [{456 12} {789 13}]]}
	}
}
