package proto

import (
	"fmt"

	"../consts"
	"../util"
)

type IdCheckReq struct {
	CommProto
	DevId    [32]byte
	Password [32]byte
}

type IdCheckRes struct {
	CommProto
	DevCode [3]byte
	CurTime [4]byte
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
	//icr.SetDevId("18002A000247373433343038")
	icr.SetDevId(consts.DevId)
	icr.SetPassword(consts.Password)
	//icr.SetPassword("czx")
	//icr.SetDevId("180030000247373433343038")
	//icr.SetPassword("cdhzkjyxgs")

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

func (icr *IdCheckReq) SetCommData(...interface{}) bool {
	icr.FillCommon(0x00)
	return true
}

func (ics *IdCheckRes) ParseInfo() bool {
	fmt.Println("## ParseInfo!")
	return true
}
