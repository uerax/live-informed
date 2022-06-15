/*
 * @Author: ww
 * @Date: 2022-06-15 07:40:42
 * @Description:
 * @FilePath: /live-informed/center/decision.go
 */
package center

import (
	"live-informed/bilibili"
	"live-informed/common"
)

type Decision struct {
	Des map[int64]int
}

var des *Decision
var Rsl = make(map[int64]string, 0)

func init() {
	des = &Decision{
		make(map[int64]int, 0),
	}
}

func (d *Decision) decision(data *bilibili.LiveInfoResp){
	tmp := make(map[int64]int, 0)
	result := make(map[int64]string, 0)
	for _, info := range data.Data {
		tmp[info.Uid] = info.LiveStatus
		if info.LiveStatus == 1 && des.Des[info.Uid] != 1{
			result[info.Uid] = common.SplicingMsg(info.Uid, info.UName)
		}
	}

	des.Des = tmp

	Rsl = result
}


func Clear() {
	Rsl = make(map[int64]string, 0)
}