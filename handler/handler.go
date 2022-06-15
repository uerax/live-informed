/*
 * @Author: ww
 * @Date: 2022-06-15 06:58:45
 * @Description:
 * @FilePath: /live-informed/handler/handler.go
 */

package handler

import (
	"live-informed/process"
	"strings"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
)

func ATMessageEventHandler(processor process.Processor) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return processor.ProcessMessage(input, data)
	}
}

