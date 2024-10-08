package cstruct

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const tagName = "cstruct"

func supportedTag(tag string) bool {
	switch tag {
	case "be", "le", "-":
		return true
	default:
		return false
	}
}

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

func toBytes(p any, isLast bool) []byte {
	var ret []byte

	s := reflect.ValueOf(p).Elem()
	st := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)
		isLastField := (i == s.NumField()-1) && isLast

		if !supportedTag(tag) || !f.CanSet() || !supportedType(ft.Type) ||
			(isDynamicType(ft.Type) && !isLastField) {
			continue
		}

		if f.Kind() == reflect.String {
			ret = append(ret, []byte(f.String())...)
			ret = append(ret, byte(0))
			continue
		}

		if f.Kind() == reflect.Slice {
			if f.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < f.Len(); j++ {
					buf := bytes.NewBuffer(make([]byte, 0, f.Type().Elem().Size()))
					buf.Write(toBytes(f.Index(j).Addr().Interface(), false))
					ret = append(ret, buf.Bytes()...)
				}
				return ret
			}

			size := int(f.Type().Elem().Size()) * f.Len()
			buf := bytes.NewBuffer(make([]byte, 0, size))

			switch tag {
			case "be":
				binary.Write(buf, binary.BigEndian, f.Interface())
			case "le":
				binary.Write(buf, binary.LittleEndian, f.Interface())
			case "-":
				binary.Write(buf, binary.NativeEndian, f.Interface())
			}

			ret = append(ret, buf.Bytes()...)

			return ret
		}

		if f.Kind() == reflect.Array {
			if f.Type().Len() == 0 {
				continue
			} else {
				switch f.Index(0).Kind() {
				case reflect.String, reflect.Slice:
					continue
				case reflect.Struct:
					if tag == "-" {
						for j := 0; j < f.Len(); j++ {
							ret = append(ret, toBytes(f.Index(j).Addr().Interface(), false)...)
						}
					}
					continue
				default:
				}
			}
		}

		buf := bytes.NewBuffer(make([]byte, 0, f.Type().Size()))

		switch tag {
		case "be":
			binary.Write(buf, binary.BigEndian, f.Interface())
		case "le":
			binary.Write(buf, binary.LittleEndian, f.Interface())
		case "-":
			if f.Kind() == reflect.Struct {
				buf.Write(toBytes(f.Addr().Interface(), i == s.NumField()-1 && isLastField))
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

func fromBytes(b []byte, p any, isLast bool, total *int) {
	s := reflect.ValueOf(p).Elem()
	st := s.Type()
	offset := 0

	defer func() {
		if total != nil {
			*total = offset
		}
	}()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ft := st.Field(i)
		tag := ft.Tag.Get(tagName)
		size := int(f.Type().Size())
		isLastField := (i == s.NumField()-1) && isLast

		if !supportedTag(tag) || !f.CanSet() || !supportedType(ft.Type) ||
			(isDynamicType(ft.Type) && !isLastField) {
			continue
		}

		if f.Kind() == reflect.String {
			f.SetString(string(b[offset:]))
			return
		}

		if f.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(f.Type(), 1, 1)
			if !supportedType(slice.Index(0).Type()) || isDynamicType(slice.Index(0).Type()) {
				return
			}
			left := len(b) - offset
			if left <= 0 {
				return
			}

			if f.Type().Elem().Kind() == reflect.Struct {
				for j := 0; offset < len(b); j++ {
					fromBytes(b[offset:], slice.Index(0).Addr().Interface(), false, &size)
					f.Set(reflect.Append(f, slice.Index(0)))
					offset += size
				}
				return
			}

			size = int(slice.Index(0).Type().Size())
			nelem := left / size
			f.Set(reflect.MakeSlice(f.Type(), nelem, nelem))
			buf := bytes.NewBuffer(b[offset:])

			switch tag {
			case "be":
				binary.Read(buf, binary.BigEndian, f.Addr().Interface())
			case "le":
				binary.Read(buf, binary.LittleEndian, f.Addr().Interface())
			case "-":
				binary.Read(buf, binary.NativeEndian, f.Addr().Interface())
			}
			return
		}

		if f.Kind() == reflect.Array {
			if f.Type().Len() == 0 {
				continue
			} else {
				fi := f.Index(0)
				if !supportedType(fi.Type()) || isDynamicType(fi.Type()) {
					continue
				}
				if fi.Kind() == reflect.Struct {
					if tag == "-" {
						for j := 0; j < f.Len(); j++ {
							fromBytes(b[offset:], f.Index(j).Addr().Interface(), false, &size)
							offset += size
						}
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

		switch tag {
		case "be":
			binary.Read(buf, binary.BigEndian, f.Addr().Interface())
		case "le":
			binary.Read(buf, binary.LittleEndian, f.Addr().Interface())
		case "-":
			if f.Kind() == reflect.Struct {
				fromBytes(buf.Bytes(), f.Addr().Interface(), i == s.NumField()-1 && isLastField, &size)
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
//
// "le": The field is serialized in little-endian byte order;
//
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
// and must not belong to another struct or array, or will be ignored.
//
// If the field type is 'struct', the tag must be "-" or will be ignored.
//
// If the field type is 'slice', the slice type is struct, all the 'slice' and 'string'
// field inside the struct will be ignored.
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
// If the field type is 'slice' or 'string', it must be the last field in the struct
// and must not belong to another struct or array, or will be ignored.
//
// If the field type is 'struct', the tag must be "-" or will be ignored.
//
// If the field type is 'slice', the slice type is struct, all the 'slice' and 'string'
// field inside the struct will be ignored.
//
// If the field type is slice or string, it must be the last field in the struct
// or will be ignored.
func FromBytes[T any](b []byte, t *T) {
	if t == nil {
		return
	}

	if reflect.TypeOf(*t).Kind() != reflect.Struct {
		return
	}

	fromBytes(b, t, true, nil)
}
