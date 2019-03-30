package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"../util"
)

const newsVer = 1

const (
	BASIC_NEWS_TYPE = iota + 0x01
	PROTECT_NEWS_TYPE
	OFFSET_NEWS_TYPE
)

type NewsReq struct {
	CommProto
	NewsBody interface{}
}

type NewsRes struct {
	CommProto
	NewsBody interface{}
}

// 基本信息
type BasicNewsReq struct {
	NewsType       uint8    // 信息类型 0x01: 基本信息；0x02：保护区信息；0x03：限位信息
	NewsVer        uint16   // 信息版本
	TowerName      [16]byte // 塔基名称
	TowerId        uint8    // 塔基ID
	TowerSwarmId   uint8    // 塔群ID
	TypePower      uint8    // 塔基类型4bit: bits[7:4], 吊绳倍率4bit: bits[3:0]
	AxesX          int16    // 坐标X
	AxesY          int16    // 坐标Y
	FrontArmLen    uint16   // 前臂长度
	BackArmLen     uint16   // 后臂长度
	ArmHeight      uint16   // 塔臂高度
	HatHeight      uint8    // 塔帽高度
	JackingAngle   uint16   // 顶升角度
	FixLongitude   float64  // 安装经度
	FixLatitude    float64  // 安装纬度
	NorthFaceAngle uint16   // 指北角度
	Version        uint8    // 版本号
	SubVer         uint8    // 子版本号
	ApkVer         uint8    // apk版本
	ApkSubVer      uint8    // apk 子版本
}

func CalTypePower(towerType uint8, power uint8) uint8 {
	var out uint8 = 0
	typeMask, _ := strconv.ParseInt("1111", 2, 8)
	towerType = towerType & uint8(typeMask)
	towerType = towerType << 4

	powerMask, _ := strconv.ParseInt("1111", 2, 8)
	power = power & uint8(powerMask)

	out = towerType | power

	return out
}

func (bnr *BasicNewsReq) SetBasicNewsReq() {
	bnr.NewsType = 0x01 // 基本信息
	bnr.NewsVer = newsVer
	util.SliceMerge(bnr.TowerName[:], 0, []byte("czx"))
	bnr.TowerId = 1
	bnr.TowerSwarmId = 2
	bnr.TypePower = CalTypePower(1, 1)
	bnr.AxesX = 10
	bnr.AxesY = 10
	bnr.FrontArmLen = 1000
	bnr.BackArmLen = 5000
	bnr.ArmHeight = 40
	bnr.HatHeight = 5
	bnr.JackingAngle = 20
	bnr.FixLongitude = 111
	bnr.FixLatitude = 222
	bnr.NorthFaceAngle = 123
	bnr.Version = newsVer
	bnr.SubVer = 0
	bnr.ApkVer = 8
	bnr.ApkSubVer = 0
}

type BasicNewsRes struct {
	NewsType uint8  // 信息类型
	NewsVer  uint16 // 信息版本号
}

// 保护区信息
type ProtectedNewsReq struct {
	NewsType   uint8    // 信息类型 0x01: 基本信息；0x02：保护区信息；0x03：限位信息
	NewsVer    uint16   // 信息版本
	ZoneNum    uint8    // 保护区个数
	ZoneTSN    uint8    // 保护区类型(T): 0 ~ 禁行区, 1 ~ 障碍物 bits[7:6]; 保护区序号(S) bits[6:4]; 保护区元素个数(N):bits[3:0]
	ZoneName   [16]byte // 保护区名称
	ZoneId     uint8    // 保护区ID
	ArchType   uint8    // 保护区建筑类别
	ZoneHeight uint16   // 保护区高度
	ElemType   uint8    // 保护区元素信息 0x00: 点; 0x01: 圆弧
	ElemSite   interface{}
}

func (pnr *ProtectedNewsReq) SetProtectedNewsReq() {
	pnr.NewsType = 0x02
	pnr.NewsVer = newsVer
	pnr.ZoneNum = 3
	pnr.ZoneTSN = 0
	util.SliceMerge(pnr.ZoneName[:], 0, []byte("czx"))
	pnr.ZoneId = 1
	pnr.ArchType = 2
	pnr.ZoneHeight = 100
	pnr.ElemType = 0x00

	pnr.ElemSite = SpotoElemSite{
		AxesX: 10,
		AxesY: 10,
	}

}

type SpotoElemSite struct {
	AxesX int16 // X坐标
	AxesY int16 // Y坐标
}

type ArcElemSite struct {
	AxesX      int16  // X坐标
	AxesY      int16  // Y坐标
	Radius     uint16 // 圆半径
	StartAngle uint16 // 起点角度
	EndAngle   uint16 // 终点角度
}

type ProtectedNewsRes struct {
	NewsType   uint8  // 信息类型
	NewsVer    uint16 // 信息版本号
	ZoneSerial uint8  // 保护区序号
}

// 限位区信息
type OffsetNewsReq struct {
	NewsType         uint8  // 信息类型 0x01: 基本信息；0x02：保护区信息；0x03：限位信息
	NewsVer          uint16 // 信息版本
	LeftOffset       int16  // 左限位
	RightOffset      int16  // 右限位
	RemoteOffset     uint8  // 远限位
	NearOffset       uint8  // 近限位
	HightOffset      uint8  // 高限位
	StartWeighOffset uint16 // 起重量限位
	MaxRangeOffset   uint16 // 最大幅度起重量限位
	ForceOffset      uint16 // 力矩限位
	SensorEnable     uint8  // 传感器使能标识 bit0: 回转传感；1：幅度；2：高度；3：称重; 4: 行走；5：风速；6：塔身倾斜
}

func (onr *OffsetNewsReq) SetOffsetNewsReq() {
	onr.NewsType = 0x03
	onr.NewsVer = newsVer
	onr.LeftOffset = 1
	onr.RightOffset = 2
	onr.RemoteOffset = 3
	onr.NearOffset = 4
	onr.HightOffset = 5
	onr.StartWeighOffset = 6
	onr.MaxRangeOffset = 7
	onr.ForceOffset = 8
	onr.SensorEnable = 9
}

type OffsetNewsRes struct {
	NewsType uint8  // 信息类型
	NewsVer  uint16 // 信息版本号
}

func (rdr *NewsReq) Initilize(args ...interface{}) bool {
	rdr.SetFrameType(0x04) // 信息传输帧, 帧类型为0x04
	return true
}

func (rdr *NewsReq) SetPrivate(args ...interface{}) bool {
	// 基本信息
	newsType, _ := args[0].(int)
	if newsType == 0x01 {
		bnr := BasicNewsReq{}
		bnr.SetBasicNewsReq()
		size := binary.Size(bnr)
		rdr.SetInfoLen(uint16(size))
		buffer := new(bytes.Buffer)
		binary.Write(buffer, binary.BigEndian, bnr)
		rdr.Info = make([]byte, size, size)
		util.SliceMerge(rdr.Info[:], 0, buffer.Bytes())
	}

	if newsType == 0x02 {
		pnr := ProtectedNewsReq{}
		pnr.SetProtectedNewsReq()
		size := util.SizeOfStruct(pnr)
		rdr.SetInfoLen(uint16(size))
		var buff = util.BytesOfStruct(pnr)
		rdr.Info = make([]byte, size, size)
		util.SliceMerge(rdr.Info[:], 0, buff)
	}

	if newsType == 0x03 {
		onr := OffsetNewsReq{}
		onr.SetOffsetNewsReq()
		size := util.SizeOfStruct(onr)
		rdr.SetInfoLen(uint16(size))
		var buff = util.BytesOfStruct(onr)
		rdr.Info = make([]byte, size, size)
		util.SliceMerge(rdr.Info[:], 0, buff)
	}

	return true
}

func (rdr *NewsReq) SetCommData(args ...interface{}) bool {
	rdr.FillCommon()
	return true
}

func (rds *NewsRes) ParseInfo() bool {
	fmt.Println("## NewsRes ParseInfo!")
	return true
}
