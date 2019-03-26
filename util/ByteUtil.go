package util

import "fmt"

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
