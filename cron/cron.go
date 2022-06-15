/*
 * @Author: ww
 * @Date: 2022-06-15 10:31:39
 * @Description:
 * @FilePath: \live-informed\cron\cron.go
 */
package cron

import (
	"live-informed/center"
	"live-informed/process"
	"time"

	"github.com/tencent-connect/botgo/log"
)

func Init() {
	go Do()
}

func Do() {
	t := time.NewTimer(10 * time.Second)
	select {
	case <-t.C:
		log.Infof("----- 定时探测决策 -----")
		center.TaskStart()
		process.Process.SendMsgs(center.Rsl)
		center.Clear()
	}
}
