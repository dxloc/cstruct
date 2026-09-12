package cstruct

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const tagName = "cstruct"

func supportedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Uint, reflect.Map,
		reflect.Uintptr, reflect.Pointer, reflect.UnsafePointer,
		reflect.Interface, reflect.Chan, reflect.Func:
		return false
	default:
		return true
	}
}

func isDynamicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Slice, reflect.String:
		return true
	default:
		return false
	}
}

func marshal(p any, endian binary.ByteOrder, isLast bool) []byte {
	var ret []byte

	s := reflect.ValueOf(p).Elem()
	st := s.Type()

	for i := range s.NumField() {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)
		isLastField := (i == s.NumField()-1) && isLast
		edn := endian

		switch tag {
		case "be":
			edn = binary.BigEndian
		case "le":
			edn = binary.LittleEndian
		case "-":
			edn = binary.NativeEndian
		}

		if !f.CanSet() || !supportedType(ft.Type) ||
			(isDynamicType(ft.Type) && !isLastField) {
			continue
		}

		switch f.Kind() {
		case reflect.String:
			ret = append(ret, []byte(f.String())...)
			ret = append(ret, byte(0))
			return ret

		case reflect.Slice:
			switch f.Type().Elem().Kind() {
			case reflect.Uint8:
				ret = append(ret, f.Bytes()...)
				return ret
			case reflect.Struct:
				for j := 0; j < f.Len(); j++ {
					buf := bytes.NewBuffer(make([]byte, 0, f.Type().Elem().Size()))
					buf.Write(marshal(f.Index(j).Addr().Interface(), edn, false))
					ret = append(ret, buf.Bytes()...)
				}
				return ret
			}

			size := int(f.Type().Elem().Size()) * f.Len()
			buf := bytes.NewBuffer(make([]byte, 0, size))
			binary.Write(buf, edn, f.Interface())

			ret = append(ret, buf.Bytes()...)
			return ret

		case reflect.Array:
			if f.Type().Len() == 0 {
				continue
			} else {
				switch f.Index(0).Kind() {
				case reflect.String, reflect.Slice:
					continue
				case reflect.Uint8:
					ret = append(ret, f.Bytes()...)
					continue
				case reflect.Struct:
					for j := range f.Len() {
						ret = append(ret, marshal(f.Index(j).Addr().Interface(), edn, false)...)
					}
					continue
				}
			}
		}

		buf := bytes.NewBuffer(make([]byte, 0, f.Type().Size()))

		if f.Kind() == reflect.Struct {
			buf.Write(marshal(f.Addr().Interface(), edn, i == s.NumField()-1 && isLastField))
		} else {
			binary.Write(buf, edn, f.Interface())
		}

		ret = append(ret, buf.Bytes()...)
	}

	return ret
}

func unmarshal(b []byte, p any, endian binary.ByteOrder, isLast bool, total *int) {
	s := reflect.ValueOf(p).Elem()
	st := s.Type()
	offset := 0

	defer func() {
		if total != nil {
			*total = offset
		}
	}()

	for i := range s.NumField() {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)
		size := int(f.Type().Size())
		isLastField := (i == s.NumField()-1) && isLast
		edn := endian

		switch tag {
		case "be":
			edn = binary.BigEndian
		case "le":
			edn = binary.LittleEndian
		case "-":
			edn = binary.NativeEndian
		}

		if !f.CanSet() || !supportedType(ft.Type) ||
			(isDynamicType(ft.Type) && !isLastField) {
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(string(b[offset:]))
			return

		case reflect.Slice:
			slice := reflect.MakeSlice(f.Type(), 1, 1)
			if !supportedType(slice.Index(0).Type()) || isDynamicType(slice.Index(0).Type()) {
				return
			}
			left := len(b) - offset
			if left <= 0 {
				return
			}

			switch f.Type().Elem().Kind() {
			case reflect.Uint8:
				f.SetBytes(b[offset:])
				return
			case reflect.Struct:
				for j := 0; offset < len(b); j++ {
					unmarshal(b[offset:], slice.Index(0).Addr().Interface(), edn, false, &size)
					f.Set(reflect.Append(f, slice.Index(0)))
					offset += size
				}
				return
			}

			size = int(slice.Index(0).Type().Size())
			nelem := left / size
			f.Set(reflect.MakeSlice(f.Type(), nelem, nelem))
			buf := bytes.NewBuffer(b[offset:])
			binary.Read(buf, edn, f.Addr().Interface())

			return

		case reflect.Array:
			if f.Type().Len() == 0 {
				continue
			} else {
				fi := f.Index(0)
				if !supportedType(fi.Type()) || isDynamicType(fi.Type()) {
					continue
				}
				switch fi.Kind() {
				case reflect.Uint8:
					reflect.Copy(f, reflect.ValueOf(b[offset:]))
					offset += f.Type().Len()
					continue
				case reflect.Struct:
					for j := range f.Len() {
						unmarshal(b[offset:], f.Index(j).Addr().Interface(), edn, false, &size)
						offset += size
					}
					continue
				}
			}
		}

		if offset+size > len(b) {
			size = len(b) - offset
		}
		if size <= 0 {
			return
		}
		buf := bytes.NewBuffer(b[offset : offset+size])

		if f.Kind() == reflect.Struct {
			unmarshal(buf.Bytes(), f.Addr().Interface(), edn, i == s.NumField()-1 && isLastField, &size)
		} else {
			binary.Read(buf, edn, f.Addr().Interface())
		}

		offset += size
	}
}

/*
Marshal takes a pointer to a struct and returns a byte slice containing the
serialized fields of the struct, according to the "cstruct" struct tags.

The "cstruct" tag can have one of the following values:

  - "be": The field is serialized in big-endian byte order;

  - "le": The field is serialized in little-endian byte order;

  - "-": The field is serialized in the native byte order of the system or follows
    the parent struct endianess.

If the tag is not set, the field is serialized in the native byte order of the
system.

Supported fixed-size types like 'int8', 'int16', 'int32', 'int64', 'uint8',
'uint16', 'uint32', 'uint64', 'float32', 'float64', 'complex64', 'complex128'.

Embedded structs and arrays are also supported.

Strings and slices are partially supported.

If the field is not exported, it will be ignored.

If the field type is 'bool', 'int', 'uint', 'map', 'pointer', 'unsafe.Pointer',
'uintptr', 'interface', 'chan' or 'func', the field is ignored.

If the field type is 'slice' or 'string', it must be the last field in the struct
and must not belong to another struct or array, or will be ignored.

If the field type is 'slice', the slice type is struct, all the 'slice' and 'string'
field inside the struct will be ignored.

The function returns nil if the input is a nil pointer or not a pointer
to a struct, or the struct cannot be converted to byte array.
*/
func Marshal[T any](t *T) []byte {
	if t == nil {
		return nil
	}

	if reflect.TypeOf(*t).Kind() != reflect.Struct {
		return nil
	}

	return marshal(t, binary.NativeEndian, true)
}

/*
Unmarshal takes a byte slice containing serialized fields of a struct and
sets the fields of the struct according to the "cstruct" struct tags.

The "cstruct" tag can have one of the following values:

  - "be": The field is serialized in big-endian byte order;

  - "le": The field is serialized in little-endian byte order;

  - "-": The field is serialized in the native byte order of the system or follows
    the parent struct endianess.

If the tag is not set, the field is serialized in the native byte order of the
system.

Supported fixed-size types like 'int8', 'int16', 'int32', 'int64', 'uint8',
'uint16', 'uint32', 'uint64', 'float32', 'float64', 'complex64', 'complex128'.

Strings and slices are partially supported.

Embedded structs and arrays are also supported.

If the field is not exported, it will be ignored.

If the field type is bool, int, uint, map, pointer, unsafe.Pointer,
uintptr, interface, chan or func, the field is ignored.

If the field type is 'slice' or 'string', it must be the last field in the struct
and must not belong to another struct or array, or will be ignored.

If the field type is 'slice', the slice type is struct, all the 'slice' and 'string'
field inside the struct will be ignored.

If the field type is slice or string, it must be the last field in the struct
or will be ignored.
*/
func Unmarshal[T any](b []byte, t *T) {
	if t == nil {
		return
	}

	if reflect.TypeOf(*t).Kind() != reflect.Struct {
		return
	}

	unmarshal(b, t, binary.NativeEndian, true, nil)
}
