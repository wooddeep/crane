package proto

/*
#include <stdio.h>
#include <stdlib.h>

unsigned short crc_calc(unsigned char* pucSendBuf, unsigned short usLen)
{
	unsigned short i, j;
	unsigned short usCrc = 0xFFFF;
	for (i = 0; i < usLen; i++) {
		usCrc ^= (unsigned short)pucSendBuf[i];
		for (j = 0; j < 8; j++)	{
           if (usCrc & 1) {
				usCrc = (usCrc >> 1) ^ 0xa001;
           	} else {
				usCrc >>= 1;
			}
		}
	}
    return usCrc;
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"../consts"
	"../util"
)

type Request interface {
	Initilize(...interface{}) bool
	SetCommData(...interface{}) bool
	SetPrivate(...interface{}) bool
	StructToBytes(frameType uint16) []byte
}

type Response interface {
	BytesToStruct(message []byte, tLen int)
	ParseInfo() bool
}

type CommProto struct {
	DleStx           uint16 // 默认值为 0xFEFB
	VersionVendor    uint8
	FrameTypeInfoLen uint16
	DevCode          []byte
	Info             []byte
	CheckCode        uint16
	DleEtx           uint16 // 默认值为 0xFEFA
}

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

func (cp *CommProto) SetDleStx() {
	cp.DleStx = 0xFEFB
}

func (cp *CommProto) SetDleEtx() {
	cp.DleStx = 0xFEFA
}

// 协议版本固定为0x04
func (cp *CommProto) SetVersion(version uint8) {
	cp.VersionVendor = cp.VersionVendor | ((version & 0x000F) << 4)
}

func (cp *CommProto) SetVendor(vendor uint8) {
	cp.VersionVendor = cp.VersionVendor | (vendor & 0x000F)
}

func (cp *CommProto) SetFrameType(frameType uint16) {
	cp.FrameTypeInfoLen = cp.FrameTypeInfoLen | ((frameType & 0x001F) << 11)
}

func (cp *CommProto) SetInfoLen(infoLen uint16) {
	cp.FrameTypeInfoLen = cp.FrameTypeInfoLen | (infoLen & 0x007FF)
}

func (cp *CommProto) SetDevCode(devCode []byte, devLen int) {
	//cp.DevCode = devCode
	cp.DevCode = make([]byte, devLen)
	cp.DevCode = devCode[:]
	fmt.Println("## SetDevCode")
	util.DumpBytes(cp.DevCode, len(cp.DevCode))
}

// 帧校验方式X16+X15+X2+1
func (cp *CommProto) SetCheckCode(cc uint16) {
	cp.CheckCode = cc
}

func (cp *CommProto) StructToBytes(frameType uint16) []byte {
	var commLen = 25
	if frameType != 0x00 {
		commLen = 25 - (16 - 3)
	}
	var length = commLen + len(cp.Info)
	var out = make([]byte, length, length)
	var relLength = 0

	// DLE STX
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, cp.DleStx)
	util.SliceMerge(out, relLength, buffer.Bytes())
	relLength += len(buffer.Bytes())
	buffer.Reset()

	// VERSION & VENDOR
	binary.Write(buffer, binary.BigEndian, cp.VersionVendor)
	util.SliceMerge(out, relLength, buffer.Bytes())
	relLength += len(buffer.Bytes())
	buffer.Reset()

	// FrameTypeInfoLen
	binary.Write(buffer, binary.BigEndian, cp.FrameTypeInfoLen)
	util.SliceMerge(out, relLength, buffer.Bytes())
	relLength += len(buffer.Bytes())
	buffer.Reset()

	// dev code
	util.SliceMerge(out, relLength, cp.DevCode[:])
	relLength += len(cp.DevCode[:])

	// info
	util.SliceMerge(out, relLength, cp.Info[:])
	relLength += len(cp.Info[:])

	var rawData = out[2 : len(out)-4]
	var checkCode = C.crc_calc(((*C.uchar)(unsafe.Pointer(&rawData[0]))), C.ushort(len(out)-6))

	// CheckCode
	binary.Write(buffer, binary.BigEndian, uint16(checkCode)) //16910
	util.SliceMerge(out, relLength, buffer.Bytes())
	relLength += len(buffer.Bytes())
	buffer.Reset()

	// DleEtx
	binary.Write(buffer, binary.BigEndian, cp.DleEtx)
	util.SliceMerge(out, relLength, buffer.Bytes())
	relLength += len(buffer.Bytes())
	buffer.Reset()
	return out
}

// fe fb 47 08 08 7b db 05 00 7b db 05 24 f3 0c 1a cc bc fe fa 回应帧
func (cp *CommProto) BytesToStruct(message []byte, tLen int) {
	if message[0] != 0xFE || message[1] != 0xFB || message[tLen-2] != 0xFE || message[tLen-1] != 0xFA {
		fmt.Println("frame format error!")
		return
	}

	var relLength = 0

	// DLE STX (长度为2字节)
	relLength += 2

	// VERSION & VENDOR (长度为1字节)
	var vv uint8 = 0
	buf := bytes.NewReader(message[relLength : relLength+1])
	binary.Read(buf, binary.BigEndian, &vv)
	cp.VersionVendor = vv
	relLength += 1

	// FrameTypeInfoLen (长度为2字节)
	var ftil uint16 = 0
	buf = bytes.NewReader(message[relLength : relLength+2])
	binary.Read(buf, binary.BigEndian, &ftil)
	cp.FrameTypeInfoLen = ftil
	relLength += 2

	// frame type
	var frameType uint8 = uint8((cp.FrameTypeInfoLen & 0xF800) >> 11)
	// info length
	var infoLen uint16 = cp.FrameTypeInfoLen & 0x007FF

	// TODO: 这个长度3 是 设备认证的返回，可能需要切换成其他的
	// 帧类型为0x02则长度为3，帧类型其他则长度为16
	// dev code (长度为3字节)
	fmt.Printf("## frameType = %02x\n", frameType)
	var devCodeLen = 16
	if frameType != 0x00 {
		devCodeLen = 3
	}

	cp.DevCode = make([]byte, devCodeLen)
	util.SliceMerge(cp.DevCode[:], 0, message[relLength:relLength+devCodeLen])
	relLength += devCodeLen

	fmt.Println("## GetRespDevCode")
	util.DumpBytes(cp.DevCode, len(cp.DevCode))

	// info (长度变化)
	cp.Info = make([]byte, infoLen, infoLen)
	util.SliceMerge(cp.Info, 0, message[relLength:relLength+int(infoLen)])
	relLength += int(infoLen)

	// checkcode (长度为2字节)
	var remoteCheckCode uint16 = 0
	buf = bytes.NewReader(message[relLength : relLength+2])
	binary.Read(buf, binary.BigEndian, &remoteCheckCode)
	var localCheckCode = C.crc_calc(((*C.uchar)(unsafe.Pointer(&message[2]))), C.ushort(tLen-6))
	if remoteCheckCode != uint16(localCheckCode) {
		fmt.Println("## check code not match!")
	}

	return
}

func (request *CommProto) FillCommon(frameType uint16) {
	request.SetVersion(4)
	request.SetVendor(9)
	request.SetFrameType(0)
	//var devCode = [16]byte{}
	//var relDevCode = []byte("373433343038")

	//icr.SetDevId("18002A000247373433343038")
	//icr.SetPassword("czx")
	//icr.SetDevId("180030000247373433343038")

	//var relDevCode = []byte{0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x2A, 0x00, 0x02, 0x47, 0x37, 0x34, 0x33, 0x34, 0x30, 0x38}
	//var relDevCode = []byte{0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x03, 0x00, 0x02, 0x47, 0x37, 0x34, 0x33, 0x34, 0x30, 0x38}
	var relDevCode = consts.DevCode
	var devLen = 16
	if frameType != 0x00 {
		relDevCode = []byte{0x4a, 0x06, 0x9e}
		//relDevCode = []byte{0x9f, 0x3f, 0x69}
		devLen = 3
	}

	/*for i, v := range relDevCode {
		devCode[len(devCode)-len(relDevCode)+i] = v
	}*/
	request.SetDevCode(relDevCode[:], devLen)
}

func (request *CommProto) ToCommMsg(frameType uint16) []byte {
	message := request.StructToBytes(frameType)
	return message
}

func Initilize(proto Request) {
	proto.Initilize()
}

func ProtoPrivate(proto Request) {
	proto.SetPrivate()
}

func ProtoCommData(proto Request) {
	proto.SetCommData()
}
