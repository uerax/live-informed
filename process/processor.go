/*
 * @Author: ww
 * @Date: 2022-06-15 07:08:11
 * @Description:
 * @FilePath: \live-informed\process\processor.go
 */

package process

import (
	"context"
	"encoding/json"
	"fmt"
	"live-informed/common"
	"live-informed/result"
	"log"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
)

type Processor struct {
	Api openapi.OpenAPI
}

var (
	channelId string
	// 消息处理器，持有 openapi 对象
	Process Processor
)

func init() {
	channelId = "7491017"
	Process = Processor{}
}

func (p Processor) SendMsgs(msg string) {
	if len(msg) == 0 {
		return
	}
	ctx := context.Background()
	toCreate := &dto.MessageToCreate{
		MsgID:   "",
		Content: msg,
	}

	p.Api.PostMessage(ctx, channelId, toCreate)
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {

	channelId = data.ChannelID
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	toCreate := &dto.MessageToCreate{
		Content: "默认回复" + message.Emoji(307),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	// 进入到私信逻辑
	if cmd.Cmd == "dm" {
		p.dmHandler(data)
		return nil
	}

	switch cmd.Cmd {
	case "/status":
		msg := ""
		if result.Rsl["211336"] {
			msg = "香香鸡腿堡正在直播"
		} else {
			msg = "香香鸡腿堡未开播"
		}
		toCreate.Content = msg
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "/add":
		toCreate.Content = "暂不支持"
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "hi":
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "time":
		toCreate.Content = genReplyContent(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "ark":
		toCreate.Ark = genReplyArk(data)
		p.sendReply(ctx, data.ChannelID, toCreate)
	case "公告":
		p.setAnnounces(ctx, data)
	case "pin":
		if data.MessageReference != nil {
			p.setPins(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	case "emoji":
		if data.MessageReference != nil {
			p.setEmoji(ctx, data.ChannelID, data.MessageReference.MessageID)
		}
	default:
	}

	return nil
}

// ProcessInlineSearch is a function to process inline search
func (p Processor) ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := p.Api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func (p Processor) dmHandler(data *dto.WSATMessageData) {
	dm, err := p.Api.CreateDirectMessage(
		context.Background(), &dto.DirectMessageToCreate{
			SourceGuildID: data.GuildID,
			RecipientID:   data.Author.ID,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}

	toCreate := &dto.MessageToCreate{
		Content: "默认私信回复",
	}
	_, err = p.Api.PostDirectMessage(
		context.Background(), dm, toCreate,
	)
	if err != nil {
		log.Println(err)
		return
	}
}

func genReplyContent(data *dto.WSATMessageData) string {
	var tpl = `你好：%s
在子频道 %s 收到消息。
收到的消息发送时时间为：%s
当前本地时间为：%s
消息来自：%s
`

	msgTime, _ := data.Timestamp.Time()
	return fmt.Sprintf(
		tpl,
		message.MentionUser(data.Author.ID),
		message.MentionChannel(data.ChannelID),
		msgTime, time.Now().Format(time.RFC3339),
		common.GetIP(),
	)
}

func genReplyArk(data *dto.WSATMessageData) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "这是 ark 的描述信息",
			},
			{
				Key:   "#PROMPT#",
				Value: "这是 ark 的摘要信息",
			},
			{
				Key: "#LIST#",
				Obj: []*dto.ArkObj{
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这里展示的是 23 号模板",
							},
						},
					},
					{
						ObjKV: []*dto.ArkObjKV{
							{
								Key:   "desc",
								Value: "这是 ark 的列表项名称",
							},
							{
								Key:   "link",
								Value: "https://www.qq.com",
							},
						},
					},
				},
			},
		},
	}
}

func (p Processor) setEmoji(ctx context.Context, channelID string, messageID string) {
	err := p.Api.CreateMessageReaction(
		ctx, channelID, messageID, dto.Emoji{
			ID:   "307",
			Type: 1,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setPins(ctx context.Context, channelID, msgID string) {
	_, err := p.Api.AddPins(ctx, channelID, msgID)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setAnnounces(ctx context.Context, data *dto.WSATMessageData) {
	if _, err := p.Api.CreateChannelAnnounces(
		ctx, data.ChannelID,
		&dto.ChannelAnnouncesToCreate{MessageID: data.ID},
	); err != nil {
		log.Println(err)
	}
}

func (p Processor) sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := p.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
	}
}
