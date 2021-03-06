package proto

import (
	"fmt"

	"../util"
)

type RealDataReq struct {
	CommProto
	Time            [4]byte // 状态采集时间（）
	Rotation        [2]byte // 回转角度值 (0.1为单位，以下同) -3276.8 ~ 3276.7
	Latitude        [2]byte // 幅度 0 ~ 6553.5
	Height          [2]byte // 高度 -3276.8
	Weight          [2]byte // 称重 0.01为单位 0 ~655.35
	Moment          [1]byte // 力矩 0-255
	Battery         [1]byte // 电池电量 0-100
	WindSpeed       [2]byte // 风速 0 ~ 6553.5
	TowerXInc       [1]byte // 塔x倾斜度 -12.8 ~ 12.7 单位0.1
	TowerYInc       [1]byte // 塔y倾斜度 -12.8 ~ 12.7 单位0.1
	OffsetAlm       [4]byte // 限位报警，详见文档
	OtherAlm        [4]byte // 其他告警，详见文档
	TowerCrashAlm   [4]byte // 塔基碰撞告警
	DmzCrashAlm     [4]byte // 禁行区碰撞告警
	BarrierCrashAlm [4]byte // 障碍物碰撞告警
	RelayOutCode    [4]byte // 继电输出编码可全为0
}

type RealDataRes struct {
	CommProto
	Time    [4]byte // 状态采集时间（）
	Command [1]byte // 设置指令
}

func (rdr *RealDataReq) Initilize(...interface{}) bool {
	rdr.SetFrameType(0x02)
	return true
}

func (rdr *RealDataReq) SetPrivate(args ...interface{}) bool {
	timeBytes := util.CurTimeToBytes()
	util.SliceMerge(rdr.Time[:], 0, timeBytes)

	fmt.Printf("# args = %v\n", args)
	arg_0, _ := args[0].(int)
	data := byte(arg_0)
	fmt.Printf("# data = %v\n", byte(data))

	rdr.SetInfoLen(42)
	rdr.Info = make([]byte, 42, 42)
	util.SliceMerge(rdr.Info[:], 0, timeBytes)

	util.SliceMerge(rdr.Info[:], 4, []byte{data, data})
	util.SliceMerge(rdr.Info[:], 6, []byte{data, data})
	util.SliceMerge(rdr.Info[:], 8, []byte{data, data})
	util.SliceMerge(rdr.Info[:], 10, []byte{data, data})

	util.SliceMerge(rdr.Info[:], 12, []byte{data})
	util.SliceMerge(rdr.Info[:], 13, []byte{data})
	util.SliceMerge(rdr.Info[:], 14, []byte{data, data})
	util.SliceMerge(rdr.Info[:], 16, []byte{data})
	util.SliceMerge(rdr.Info[:], 17, []byte{data})

	util.SliceMerge(rdr.Info[:], 18, []byte{data, data, data, data})
	util.SliceMerge(rdr.Info[:], 22, []byte{data, data, data, data})
	util.SliceMerge(rdr.Info[:], 26, []byte{data, data, data, data})
	util.SliceMerge(rdr.Info[:], 30, []byte{data, data, data, data})
	util.SliceMerge(rdr.Info[:], 34, []byte{data, data, data, data})
	util.SliceMerge(rdr.Info[:], 38, []byte{data, data, data, data})
	//rdr.Info = buff

	return true
}

func (rdr *RealDataReq) SetCommData(...interface{}) bool {
	rdr.FillCommon(0x02)
	return true
}

func (rds *RealDataRes) ParseInfo() bool {
	fmt.Println("## RealDataRes ParseInfo!")
	return true
}
