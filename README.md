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

## Installation

To install the `cstruct` library, run the following command:

```bash
go get -u github.com/dxloc/cstruct
```

## Usage

To use the `cstruct` library, import it in your Go program and use the `Marshal()`
and `Unmarshal()` functions to serialize and deserialize your structs.

Examples is provided in the `cmd/` directory.

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

The cstruct library is released under the [The Unlicense](https://unlicense.org).