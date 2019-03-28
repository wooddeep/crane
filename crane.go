// 反射调用方法, 函数： https://www.cnblogs.com/52php/p/6337420.html
package main

import (
	"fmt"
	"net"
	"os"

	"./proto"
	"./util"
)

type StateMatrix struct {
	Request  proto.Request // 声明了interface类型
	Response proto.Response
}

var statMatrix []StateMatrix = []StateMatrix{
	{
		Request: &proto.IdCheckReq{ // 实现Protocol接口时, receiver为指针, 故而此处为指针
			CommProto: proto.CommProto{
				DleStx: 0xFEFB,
				DleEtx: 0xFEFA,
			},
		},
		Response: &proto.IdCheckRes{
			CommProto: proto.CommProto{
				DleStx: 0xFEFB,
				DleEtx: 0xFEFA,
			},
		},
	},

	{
		Request: &proto.RealDataReq{
			CommProto: proto.CommProto{
				DleStx: 0xFEFB,
				DleEtx: 0xFEFA,
			},
		},
		Response: &proto.RealDataRes{
			CommProto: proto.CommProto{
				DleStx: 0xFEFB,
				DleEtx: 0xFEFA,
			},
		},
	},
}

func FrameWork(socket *net.UDPConn) {
	for _, v := range statMatrix {
		var request = v.Request
		request.Initilize()   // 初始化请求
		request.SetPrivate()  // 框架调用每种协议的私有方法, 设置私有数据
		request.SetCommData() // 框架调用每种协议的公有方法, 设置共有数据
		message := request.StructToBytes()

		socket.Write(message)
		data := make([]byte, 1024)
		len, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("error recv data")
			return
		}

		util.DumpBytes(message)   // 请求
		util.DumpBytes(data, len) // 返回

		var response = v.Response
		response.BytesToStruct(data, len) // 填充返回信息
		response.ParseInfo()
	}
}

func main() {

	addr, err := net.ResolveUDPAddr("udp", "101.207.139.194:9999")
	if err != nil {
		fmt.Println("net.ResolveUDPAddr fail.", err)
		os.Exit(1)
	}

	socket, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("net.DialUDP fail.", err)
		os.Exit(1)
	}

	defer socket.Close()
	FrameWork(socket)
}
