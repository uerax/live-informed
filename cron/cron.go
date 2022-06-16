/*
 * @Author: ww
 * @Date: 2022-06-15 10:31:39
 * @Description:
 * @FilePath: /live-informed/cron/cron.go
 */
package cron

import (
	"live-informed/center"
	"time"

	"github.com/tencent-connect/botgo/log"
)

func Init() {
	Do()
}

func Do() {
	ticker := time.NewTicker(time.Second * 60)

	go func() {
		for {
			select {
			case <- ticker.C:
				log.Infof("----- 定时探测决策 -----")
				center.TaskStart()
			}
		}
	}()
	
}
