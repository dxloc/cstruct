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

func toBytes(p any, isLast bool) []byte {
	var ret []byte

	s := reflect.ValueOf(p).Elem()
	st := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)

		if tag == "" || (f.Kind() == reflect.Struct && tag != "-") ||
			((f.Kind() == reflect.Slice || f.Kind() == reflect.String) &&
				(i != s.NumField()-1 || !isLast)) ||
			!f.CanInterface() || !supportedType(ft.Type) {
			continue
		}

		if f.Kind() == reflect.String {
			ret = append(ret, []byte(f.String())...)
			continue
		}

		if f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				b := make([]byte, 0, f.Type().Elem().Size())
				buf := bytes.NewBuffer(b)

				switch tag {
				case "be":
					binary.Write(buf, binary.BigEndian, f.Index(j).Interface())
				case "le":
					binary.Write(buf, binary.LittleEndian, f.Index(j).Interface())
				case "-":
					if f.Index(j).Kind() == reflect.Struct {
						buf.Write(toBytes(f.Index(j).Addr().Interface(), true))
					} else {
						binary.Write(buf, binary.NativeEndian, f.Index(j).Interface())
					}
				}

				ret = append(ret, buf.Bytes()...)
			}
			return ret
		}

		b := make([]byte, 0, f.Type().Size())
		buf := bytes.NewBuffer(b)

		switch tag {
		case "be":
			binary.Write(buf, binary.BigEndian, f.Interface())
		case "le":
			binary.Write(buf, binary.LittleEndian, f.Interface())
		case "-":
			if f.Kind() == reflect.Struct {
				buf.Write(toBytes(f.Addr().Interface(), i == s.NumField()-1))
			} else {
				binary.Write(buf, binary.NativeEndian, f.Interface())
			}
		default:
			continue
		}

		ret = append(ret, buf.Bytes()...)
	}

	return ret
}

func fromBytes(b []byte, p any, isLast bool) {
	s := reflect.ValueOf(p).Elem()
	st := s.Type()
	offset := 0

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)
		size := int(f.Type().Size())

		if tag == "" || (f.Kind() == reflect.Struct && tag != "-") ||
			((f.Kind() == reflect.Slice || f.Kind() == reflect.String) &&
				(i != s.NumField()-1 || !isLast)) ||
			!f.CanSet() || !supportedType(ft.Type) {
			continue
		}

		if f.Kind() == reflect.String {
			f.SetString(string(b[offset : offset+size]))
			offset += size
			continue
		}

		if f.Kind() == reflect.Slice {
			size = len(b) - offset
			nelem := size / int(f.Type().Elem().Size())
			size = int(f.Type().Elem().Size())

			f.Set(reflect.MakeSlice(f.Type(), nelem, nelem))

			for j := 0; j < nelem; j++ {
				buf := bytes.NewBuffer(b[offset : offset+size])

				switch tag {
				case "be":
					binary.Read(buf, binary.BigEndian, f.Index(j).Addr().Interface())
				case "le":
					binary.Read(buf, binary.LittleEndian, f.Index(j).Addr().Interface())
				case "-":
					if f.Index(j).Kind() == reflect.Struct {
						fromBytes(buf.Bytes(), f.Index(j).Addr().Interface(), true)
					} else {
						binary.Read(buf, binary.NativeEndian, f.Index(j).Addr().Interface())
					}
				}

				offset += size
			}

			return
		}

		buf := bytes.NewBuffer(b[offset : offset+size])

		switch tag {
		case "be":
			binary.Read(buf, binary.BigEndian, f.Addr().Interface())
		case "le":
			binary.Read(buf, binary.LittleEndian, f.Addr().Interface())
		case "-":
			if f.Kind() == reflect.Struct {
				fromBytes(buf.Bytes(), f.Addr().Interface(), i == s.NumField()-1)
			} else {
				binary.Read(buf, binary.NativeEndian, f.Addr().Interface())
			}
		default:
			continue
		}

		offset += size
	}
}

// ToBytes takes a pointer to a struct and returns a byte slice containing the
// serialized fields of the struct, according to the "cstruct" struct tags.
//
// The "cstruct" tag can have one of the following values:
//
// "be": The field is serialized in big-endian byte order;
// "le": The field is serialized in little-endian byte order;
// "-": The field is serialized in the native byte order of the system.
//
// If the tag is not set, the field is ignored.
//
// If the field is not exported, it will be ignored.
//
// If the field type is 'bool', 'int', 'uint', 'map', 'pointer', 'unsafe.Pointer',
// 'uintptr', 'interface', 'chan' or 'func', the field is ignored.
//
// If the field type is 'slice' or 'string', it must be the last field in the struct
// or will be ignored.
//
// If the field type is 'struct', the tag must be "-" or will be ignored.
//
// The function returns nil if the input is a nil pointer or not a pointer
// to a struct, or the struct cannot be converted to byte array.
func ToBytes[T any](t *T) []byte {
	if t == nil {
		return nil
	}

	if reflect.TypeOf(*t).Kind() != reflect.Struct {
		return nil
	}

	return toBytes(t, true)
}

// FromBytes takes a byte slice containing serialized fields of a struct and
// sets the fields of the struct according to the "cstruct" struct tags.
//
// The "cstruct" tag can have one of the following values:
//
// "be": The field is serialized in big-endian byte order;
// "le": The field is serialized in little-endian byte order;
// "-": The field is serialized in the native byte order of the system.
//
// If the tag is not set, the field is ignored.
//
// If the field is not exported, it will be ignored.
//
// If the field type is bool, int, uint, map, pointer, unsafe.Pointer,
// uintptr, interface, chan or func, the field is ignored.
//
// If the field type is slice or string, it must be the last field in the struct
// or will be ignored.
//
// If the field type is struct, the tag must be "-" or will be ignored.
func FromBytes[T any](b []byte, t *T) {
	if t == nil {
		return
	}

	if reflect.TypeOf(*t).Kind() != reflect.Struct {
		return
	}

	fromBytes(b, t, true)
}
