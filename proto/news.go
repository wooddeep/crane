package proto

import (
	"fmt"
)

type NewsReq struct {
	CommProto
	NewsType uint8 // 信息类型 0x01: 基本信息；0x02：保护区信息；0x03：限位信息
	NewsBody interface{}
}

type NewsRes struct {
	CommProto
	NewsType uint8  // 信息类型
	NewsVer  uint16 // 信息版本号
	NewsBody interface{}
}

// 基本信息
type BasicNewsReq struct {
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

type BasicNewsRes struct {
}

// 保护区信息
type ProtectedNewsReq struct {
	NewsVer    uint16      // 信息版本
	ZoneNum    uint8       // 保护区个数
	ZoneTSN    uint8       // 保护区类型(T): 0 ~ 禁行区, 1 ~ 障碍物 bits[7:6]; 保护区序号(S) bits[6:4]; 保护区元素个数(N):bits[3:0]
	ZoneName   [16]byte    // 保护区名称
	ZoneId     uint8       // 保护区ID
	ArchType   uint8       // 保护区建筑类别
	ZoneHeight uint16      // 保护区高度
	ElemType   uint8       // 保护区元素信息 0x00: 点; 0x01: 圆弧
	ElemSite   interface{} // 保护区位置信息 SpotoElemSite 或者 ArcElemSite
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
	ZoneSerial uint8 // 保护区序号
}

// 限位区信息
type OffsetNewsReq struct {
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

type OffsetNewsRes struct {
}

func (rdr *NewsReq) Initilize(...interface{}) bool {
	rdr.SetFrameType(0x04) // 信息传输帧, 帧类型为0x04
	return true
}

func (rdr *NewsReq) SetPrivate(...interface{}) bool {
	return true
}

func (rdr *NewsReq) SetCommData(...interface{}) bool {
	rdr.FillCommon()
	return true
}

func (rds *NewsRes) ParseInfo() bool {
	fmt.Println("## RealDataRes ParseInfo!")
	return true
}
