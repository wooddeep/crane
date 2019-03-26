package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

// description
// 当前时间，用于同步。 格式说明如下，后续的时间均以此为标准:
// [31:26] 年，6bits，0-63(以 2010 为基数)
// [25:22] 月，4bits，1-12
// [21:17] 日，5bits，1-31
// [16:12] 小时，5bits，0-23
// [11:6] 分，6bits，0-59
// [5:0] 秒，6bits，0-59

func CurTimeToBytes() []byte {
	var out uint32 = 0

	year := uint32(time.Now().Year() - 2010)         //年
	yearMask, _ := strconv.ParseInt("111111", 2, 32) // 2进制数111111转换为32位的整型数
	year = year & uint32(yearMask)
	year = year << 26
	out = out | year

	month := uint32(time.Now().Month()) //月
	monthMask, _ := strconv.ParseInt("1111", 2, 32)
	month = month & uint32(monthMask)
	month = month << 22
	out = out | month

	day := uint32(time.Now().Day()) //日
	dayMask, _ := strconv.ParseInt("11111", 2, 32)
	day = day & uint32(dayMask)
	day = day << 17
	out = out | day

	hour := uint32(time.Now().Hour()) //小时
	hourMask, _ := strconv.ParseInt("11111", 2, 32)
	hour = hour & uint32(hourMask)
	hour = hour << 12
	out = out | hour

	minute := uint32(time.Now().Minute()) //分钟
	minuteMask, _ := strconv.ParseInt("111111", 2, 32)
	minute = minute & uint32(minuteMask)
	minute = minute << 6
	out = out | minute

	second := uint32(time.Now().Second()) //秒
	secondMask, _ := strconv.ParseInt("111111", 2, 32)
	second = second & uint32(secondMask)
	out = out | second

	fmt.Printf("%032b \n", out)

	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, out)

	return buffer.Bytes()
}
