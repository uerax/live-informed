/*
 * @Author: ww
 * @Date: 2022-06-15 08:27:32
 * @Description:
 * @FilePath: /live-informed/center/api.go
 */
package center

import "live-informed/result"


func TaskStart() {
	tasks.Detection()
}

func GetStatus() bool {
	return result.Rsl["211336"]
}