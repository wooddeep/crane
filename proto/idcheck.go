package proto

import (
	"fmt"

	"../util"
)

type IdCheckReq struct {
	CommProto
	DevId    [32]byte
	Password [32]byte
}

type IdCheckRes struct {
	CheckRet byte
	DevCode  [3]byte
	CurTime  [4]byte
}

func (icr *IdCheckReq) SetDevId(devId string) {
	var idInfo = util.StrToBytes(devId, 32) // hdgd
	util.SliceMerge(icr.DevId[:], 0, idInfo)
}

func (icr *IdCheckReq) SetPassword(password string) {
	var passInfo = util.StrToBytes(password, 32) // czx
	util.SliceMerge(icr.Password[:], 0, passInfo)
}

func (icr *IdCheckReq) Initilize(...interface{}) bool {
	icr.SetDevId("18002A000247373433343038")
	icr.SetPassword("czx")
	return true
}

func (icr *IdCheckReq) SetPrivate(...interface{}) bool {
	accLen := len(icr.DevId)
	len := len(icr.DevId) + len(icr.Password)
	buff := make([]byte, len, len)

	fmt.Println(len)
	util.SliceMerge(buff, 0, icr.DevId[:])
	util.SliceMerge(buff, accLen, icr.Password[:])
	icr.SetInfoLen(uint16(len))
	icr.Info = buff
	return true
}