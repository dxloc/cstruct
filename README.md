# cstruct

A Go library for serializing and deserializing structs to and from byte slices.

## Overview

The `cstruct` library provides a simple way to serialize and deserialize structs
to and from byte slices. It uses the `reflect` package to inspect the struct fields
and determine how to serialize and deserialize them.

## Features

* Serialize structs to byte slices
* Deserialize byte slices to structs
* Supports various data types, including integers, floats, strings, and slices
* Allows for custom serialization and deserialization of struct fields using tags

## Usage

To use the `cstruct` library, import it in your Go program and use the `ToBytes()`
and `FromBytes()` functions to serialize and deserialize your structs.

```go
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
	fmt.Println(b) // [200 1 0 0 123 0 0 0 0 0 0 1 0 0 0 2 0 0 0 3 0 0 0 4 72 101 108 108 111 44 32 87 111 114 108 100 33 0]

	var c MyStruct2
	cstruct.FromBytes(b, &c)

	fmt.Println(c) // {456 {123 [1 2 3 4] Hello, World!}}
}
```

## Tags

The cstruct library uses tags to determine how to serialize and deserialize struct
fields. The following tags are supported:

* "be": Serialize the field in big-endian byte order
* "le": Serialize the field in little-endian byte order
* "-": Serialize the field in the native byte order of the system

You can add these tags to your struct fields to customize their serialization and
deserialization.

## Contributing

Contributions to the cstruct library are welcome. If you have a bug fix or feature
request, please open an issue or submit a pull request.

## License

The cstruct library is licensed under the MIT License. See the LICENSE file for details.