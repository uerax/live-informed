/*
 * @Author: ww
 * @Date: 2022-06-15 01:27:09
 * @Description:
 * @FilePath: /live-informed/main.go
 */
package main

import (
	"context"
	"log"
	"time"

	"live-informed/config"
	"live-informed/cron"
	"live-informed/handler"
	"live-informed/process"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

func Init() {
	config.Init()
	cron.Init()
}

func main() {
	Init()
	ctx := context.Background()
	
	botToken := token.BotToken(config.GetAppId(), config.GetToken())
	
	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)

	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}

	 process.Process = process.Processor{Api: api}

	// websocket.RegisterResumeSignal(syscall.SIGUSR1)
	// 根据不同的回调，生成 intents
	intent := websocket.RegisterHandlers(
		// at 机器人事件，目前是在这个事件处理中有逻辑，会回消息，其他的回调处理都只把数据打印出来，不做任何处理
		handler.ATMessageEventHandler(process.Process),
		// 如果想要捕获到连接成功的事件，可以实现这个回调
		//ReadyHandler(),
		// 连接关闭回调
		//ErrorNotifyHandler(),
		// 频道事件
		//GuildEventHandler(),
		// 成员事件
		//MemberEventHandler(),
		// 子频道事件
		//ChannelEventHandler(),
		// 私信，目前只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		//DirectMessageHandler(),
		// 频道消息，只有私域才能够收到这个，如果你的机器人不是私域机器人，会导致连接报错，那么启动 example 就需要注释掉这个回调
		//CreateMessageHandler(),
		// 互动事件
		//InteractionHandler(),
		// 发帖事件
		//ThreadEventHandler(),
	)
	// 指定需要启动的分片数为 2 的话可以手动修改 wsInfo
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}


}