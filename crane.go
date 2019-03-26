// 反射调用方法, 函数： https://www.cnblogs.com/52php/p/6337420.html
package main

import (
	"fmt"
	"net"
	"os"

	"./proto"
	"./util"
)

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

	var idVerify = proto.IdCheckReq{
		CommProto: proto.CommProto{
			DleStx: 0xFEFB,
			DleEtx: 0xFEFA,
		},
	}

	proto.Initilize(&idVerify)
	proto.ProtoEntry(&idVerify) // 框架调用每种协议的私有方法

	idVerify.FillCommon()           // 设置通用信息
	message := idVerify.ToCommMsg() // 生成通信用的数据

	socket.Write(message)
	data := make([]byte, 1024)
	len, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("error recv data")
		return
	}

	util.DumpBytes(message)        // 请求
	util.DumpBytesByLen(data, len) // 返回

	// data
	var idVerifyResp = proto.CommProto{
		DleStx: 0xFEFB,
		DleEtx: 0xFEFA,
	}
	// fe fb 47 08 08 7b db 05 00 7b db 05 24 f3 0c 1a cc bc fe fa 回应帧

	idVerifyResp.BytesToStruct(data, len)

	fmt.Printf("from %s:%s\n", remoteAddr.String(), string(data))
}
