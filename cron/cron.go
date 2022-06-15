/*
 * @Author: ww
 * @Date: 2022-06-15 10:31:39
 * @Description:
 * @FilePath: /live-informed/cron/cron.go
 */
package cron

import (
	"live-informed/center"
	"live-informed/process"
	"time"
)

func Init() {
	go Do()
}

func Do() {
	timeChannel := time.NewTimer(10 * time.Second)
  	select {
    case <-timeChannel.C:
       go do()
  }
}

func do() {
	center.TaskStart()
	process.Process.SendMsgs(center.Rsl)
	center.Rsl = make(map[int64]string, 0)
}