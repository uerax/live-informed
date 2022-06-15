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
	List map[int64]struct{}
}

var tasks *Task

func init() {
	tasks = &Task{
		make(map[int64]struct{}),
	}
	tasks.List[211336] = struct{}{}
}

func (t *Task) Detection() {

	if len(t.List) == 0 {
		return
	}

	tmp := make([]int64, 0, len(t.List))

	for k, _ := range t.List {
		log.Infof("detect %d", k)
		tmp = append(tmp, k)
	}

	res, err := bilibili.GetLiveInfo(tmp)
	if err != nil {
		return
	}

	des.decision(res)

}

func (t *Task) AddTask(uid int64) {
	t.List[uid] = struct{}{}
}

func (t *Task) DelTask(uid int64) {
	delete(t.List, uid)
}
