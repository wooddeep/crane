package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type Error struct {
	ErrCode int
	ErrMsg  string
}

func NewError(code int, msg string) *Error {
	return &Error{ErrCode: code, ErrMsg: msg}
}

func (err *Error) Error() string {
	return err.ErrMsg
}

func StrToBytes(content string, length int) []byte {
	raw := []byte(content)
	var target = make([]byte, length, length)
	for i, v := range raw {
		target[i] = v
	}
	return target
}

func DumpBytes(args ...interface{}) {
	var bytes []byte
	var length int
	var ok bool
	if len(args) >= 1 {
		bytes, ok = args[0].([]byte)
		if ok != true {
			fmt.Println("## [0] parameter error!")
			return
		}
		length = len(bytes)
	}

	if len(args) >= 2 {
		length, ok = args[1].(int)
		if ok != true {
			fmt.Println("## [1] parameter error!")
			return
		}
	}

	fmt.Println("-----------")
	for i, v := range bytes {
		if i >= length {
			break
		}
		fmt.Printf("%02x ", v)
	}
	fmt.Println("\n-----------")
}

func SliceMerge(pSlice []byte, pIndex int, cSlice []byte) error {
	if len(pSlice)-pIndex+1 < len(cSlice) {
		return NewError(1, "length error!")
	}

	for i, v := range cSlice {
		pSlice[i+pIndex] = v
	}
	return nil
}

func SizeOf(stru interface{}) {
	v := reflect.ValueOf(stru)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			fmt.Println(f.String())
		case reflect.Int:
			fmt.Println(f.Int())
		}
	}
}

func SizeOfPlain(kind interface{}) int {
	var size = 0
	switch kind {
	case reflect.Int:
		size = 4
	case reflect.Int8:
		size = 1
	case reflect.Int16:
		size = 2
	case reflect.Int32:
		size = 4
	case reflect.Int64:
		size = 8
	case reflect.Uint:
		size = 4
	case reflect.Uint8:
		size = 1
	case reflect.Uint16:
		size = 2
	case reflect.Uint32:
		size = 4
	case reflect.Uint64:
		size = 8
	case reflect.Float32:
		size = 4
	case reflect.Float64:
		size = 8
	}
	return size
}

// https://stackoverflow.com/questions/42151307/how-to-determine-the-element-type-of-slice-interface
// https://stackoverflow.com/questions/52044982/get-the-length-of-a-slice-of-unknown-type
func SizeOfStruct(stru interface{}) int {
	var size = 0
	v := reflect.ValueOf(stru)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		k := f.Kind()

		if k <= reflect.Complex128 {
			size += SizeOfPlain(k)
		}

		if k == reflect.Interface {
			size += SizeOfStruct(f.Interface())
		}

		if k == reflect.Struct {
			size += SizeOfStruct(f.Interface())
		}

		if k == reflect.Array {
			len := reflect.ValueOf(f.Interface()).Type().Len()  // 数组长度
			kind := reflect.TypeOf(f.Interface()).Elem().Kind() // 数组原因类型
			if kind <= reflect.Complex128 {
				size += len * SizeOfPlain(kind)
			} else {
				size += len * SizeOfStruct(f.Interface())
			}
		}
	}
	return size
}

// var kind = reflect.TypeOf(value).Kind()
func BytesOfPlain(value interface{}) []byte {
	var out = make([]byte, 0)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, value)
	out = append(out, buffer.Bytes()...)

	return out
}

// https://segmentfault.com/q/1010000000198391

func BytesOfStruct(stru interface{}) []byte {
	var out = make([]byte, 0)
	v := reflect.ValueOf(stru)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)

		k := f.Kind()
		if k <= reflect.Complex128 {
			buffer := new(bytes.Buffer)
			binary.Write(buffer, binary.BigEndian, f.Interface())
			out = append(out, buffer.Bytes()...)
		}

		if k == reflect.Interface {
			out = append(out, BytesOfStruct(f.Interface())...)
		}

		if k == reflect.Struct {
			out = append(out, BytesOfStruct(f.Interface())...)
		}

		if k == reflect.Array {
			len := reflect.ValueOf(f.Interface()).Type().Len()
			kind := reflect.TypeOf(f.Interface()).Elem().Kind()
			for i := 0; i < len; i++ {
				iv := reflect.ValueOf(f.Interface())
				if kind <= reflect.Complex128 {
					out = append(out, BytesOfPlain(iv.Index(i).Interface())...) // 根据下标获取数组的元素值

				} else {
					out = append(out, BytesOfStruct(iv.Index(i).Interface())...)

				}
			}
		}

	}
	return out
}
