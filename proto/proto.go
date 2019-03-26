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

	"../util"
)

type Protocol interface {
	Initilize(...interface{}) bool
	SetPrivate(...interface{}) bool
}

type CommProto struct {
	DleStx           uint16 // 默认值为 0xFEFB
	VersionVendor    uint8
	FrameTypeInfoLen uint16
	DevCode          [16]byte
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

func (cp *CommProto) SetDevCode(devCode [16]byte) {
	cp.DevCode = devCode
}

// 帧校验方式X16+X15+X2+1
func (cp *CommProto) SetCheckCode(cc uint16) {
	cp.CheckCode = cc
}

func (cp *CommProto) StructToBytes() []byte {
	var length = 25 + len(cp.Info)
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
	util.DumpBytes(rawData)
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

	util.DumpBytes(out)

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
	_ = frameType
	// info length
	var infoLen uint16 = cp.FrameTypeInfoLen & 0x007FF

	// dev code (长度为3字节)
	util.SliceMerge(cp.DevCode[:], 0, message[relLength:relLength+3]) // TODO: 这个长度3 是 设备认证的返回，可能需要切换成其他的
	relLength += 3

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

func (request *CommProto) FillCommon() {
	request.SetVersion(4)
	request.SetVendor(7)
	request.SetFrameType(0)
	var devCode = [16]byte{}
	var relDevCode = []byte("373433343038")
	//var relDevCode = []byte{0x19, 0x00, 0x21, 0x00, 0x02, 0x47, 0x37, 0x34, 0x33, 0x34, 0x30, 0x38}
	//var relDevCode = []byte{0x19, 0x00, 0x21, 0x00, 0x02, 0x47, 0x37, 0x34, 0x33, 0x34, 0x30, 0x38}
	for i, v := range relDevCode {
		devCode[len(devCode)-len(relDevCode)+i] = v
	}
	request.SetDevCode(devCode)
}

func (request *CommProto) ToCommMsg() []byte {
	message := request.StructToBytes()
	return message
}

func Initilize(proto Protocol) {
	proto.Initilize()
}

func ProtoEntry(proto Protocol) {
	proto.SetPrivate()
}
