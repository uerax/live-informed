/*
 * @Author: ww
 * @Date: 2022-06-15 07:10:31
 * @Description:
 * @FilePath: /live-informed/common/common.go
 */
package common

import (
	"fmt"
	"net"
)

const (
	msg = "@所有人 %s 已开播 https://live.bilibili.com/%d"
)

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.1"
}

func SplicingMsg(roomId int64, uName string) string {
	return fmt.Sprintf(msg, uName, roomId)
}