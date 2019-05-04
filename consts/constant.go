package consts

import (
	"strconv"

	"../util"
)

//const DevId = "180030000247373433343038"
//const Password = "cdhzkjyxgs"

const DevId = "18002A000247373433343038"
const Password = "cdczx123"

var DevCode = []byte{}

func init() {
	// var relDevCode = []byte{0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x2A, 0x00, 0x02, 0x47, 0x37, 0x34, 0x33, 0x34, 0x30, 0x38}
	DevCode = make([]byte, 16)
	_ = DevCode
	len := len(DevId)
	j := 4
	for i := 0; i < len; i += 2 {
		value, _ := strconv.ParseInt(DevId[i:i+2], 16, 8)
		DevCode[j] = byte(value)
		j += 1
	}

	util.DumpBytes(DevCode)
}
