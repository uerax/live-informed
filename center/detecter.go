/*
 * @Author: ww
 * @Date: 2022-06-15 07:24:27
 * @Description:
 * @FilePath: \live-informed\center\detecter.go
 */
package center

import (
	"live-informed/bilibili"

	"github.com/tencent-connect/botgo/log"
)

// 探测列表
type Task struct {
	List string
}

var tasks *Task

func init() {
	tasks = &Task{
		List: "211336",
	}
}

func (t *Task) Detection() {

	if t.List == "" {
		return
	}

	log.Infof("detect %s", t.List)

	isLiving, err := bilibili.UserIsLiving(t.List)
	if err != nil {
		return
	}

	decision(isLiving, t.List)

}
