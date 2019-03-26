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

func DumpBytes(bytes []byte) {
	fmt.Println("-----------")
	for _, v := range bytes {
		fmt.Printf("%02x ", v)
	}
	fmt.Println("\n-----------")
}

func DumpBytesByLen(bytes []byte, len int) {
	fmt.Println("-----------")
	for i, v := range bytes {
		if i >= len {
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
