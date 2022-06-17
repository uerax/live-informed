/*
 * @Author: ww
 * @Date: 2022-06-15 07:40:42
 * @Description:
 * @FilePath: \live-informed\center\decision.go
 */
package center

import (
	"live-informed/process"

	"github.com/tencent-connect/botgo/log"
)

var Rsl = make(map[string]bool, 0)

func decision(isLiving bool, mid string) {
	log.Infof("%s 直播状态探测结果为：%t", mid, isLiving)
	if r, ok := Rsl[mid]; ok {
		if !r && isLiving {
			process.Process.SendMsgs("香香鸡腿堡 已开播")
		}
		if r && !isLiving {
			process.Process.SendMsgs("香香鸡腿堡 已下播")
		}
	}
	Rsl[mid] = isLiving
}
