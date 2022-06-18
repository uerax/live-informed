/*
 * @Author: ww
 * @Date: 2022-06-15 07:40:42
 * @Description:
 * @FilePath: /live-informed/center/decision.go
 */
package center

import (
	"live-informed/process"
	"live-informed/result"
)

func decision(isLiving bool, mid string) {
	if r, ok := result.Rsl[mid]; ok {
		if !r && isLiving {
			process.Process.SendMsgs("香香鸡腿堡 已开播")
		}
		if r && !isLiving {
			process.Process.SendMsgs("香香鸡腿堡 已下播")
		}
	}
	result.Rsl[mid] = isLiving
}