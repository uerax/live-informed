/*
 * @Author: ww
 * @Date: 2022-06-15 07:40:42
 * @Description:
 * @FilePath: /live-informed/center/decision.go
 */
package center

import "live-informed/process"

var Rsl = make(map[string]bool, 0)

func decision(isLiving bool, mid string) {
	if r, ok := Rsl[mid]; ok {
		if !r && isLiving {
			process.Process.SendMsgs("香香鸡腿堡已开播")
		}
	}
	Rsl[mid] = isLiving
}