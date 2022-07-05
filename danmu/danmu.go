/*
 * @Author: UerAx
 * @Date: 2022-07-05 18:54:12
 * @FilePath: \live-informed\danmu\danmu.go
 * Copyright (c) 2022 by UerAx uerax@live.com, All Rights Reserved.
 */
package danmu

import (
	"fmt"
	"live-informed/config"
	"live-informed/process"
	"live-informed/sdk/danmu-play/danmu"
	"live-informed/sdk/danmu-play/model"
)

func Init() {
	dm := danmu.NewBiliRoom(config.GetRoomId())
	dm.MsgHandler = func(mi *model.MessageInfo) error {
		process.Process.SendMsgs(fmt.Sprintf("[弹幕] %s: %s", mi.Info.([]interface{})[2].([]interface{})[1], mi.Info.([]interface{})[1]))
		return nil
	}
	go dm.Start()
	go dm.DanmuHandler()
}
