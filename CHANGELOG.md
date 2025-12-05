# Change log

## cstruct v0.1.1

### Bug fixes

* `panic: runtime error: slice bounds out of range` when the size of string or slice is
less than the actual size of the data type

* Add '\0' to the end of the string when converting

## cstruct v0.1.2

* fix slice of struct bugs

## cstruct v0.1.3

* fix array bugs

## cstruct v0.1.4

* fix slice type size

## cstruct v0.1.5

* Optimize deserializing []byte

## cstruct v0.1.6

* Optimize []byte SER/DES

## cstruct v0.1.7

* Fix bytes array bugs

## cstruct v0.1.8

* Fix bytes array bugs in cstruct.FromBytes() function

## cstruct v1.0.0

* Rename functions to cstruct.Marshal() and cstruct.Unmarshal()

## cstruct v1.1.0

* Serialize and deserialize with system native endian when the tag is not set
* Use switch case instead of if-else for type detection
* Fix README license part

## cstruct v1.1.1

* Fix tag parsing of slice field
