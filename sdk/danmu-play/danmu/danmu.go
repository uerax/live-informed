/*
 * @Author: ww
 * @Date: 2022-07-02 02:25:38
 * @Description:
 * @FilePath: \live-informed\sdk\danmu-play\danmu\danmu.go
 */
package danmu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"live-informed/sdk/danmu-play/global"
	"live-informed/sdk/danmu-play/model"
	"live-informed/sdk/danmu-play/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var ulog = global.Log

type BiliRoom struct {
	roomId     string
	realId     string
	address    string
	token      string
	conn       *websocket.Conn
	recMsg     chan []byte
	OutMsg     chan []byte
	isAlive    bool
	timeout    int
	MsgHandler func(mi *model.MessageInfo) error
}

func (t *BiliRoom) New(roomId string) {
	t = NewBiliRoom(roomId)
}

func NewBiliRoom(roomId string) *BiliRoom {
	var recMsg = make(chan []byte, 10)
	var OutMsg = make(chan []byte, 10)
	var conn *websocket.Conn
	var realId string
	if len(roomId) < 5 {
		realId, _ = getRealId(roomId)
		ulog.Infof("房间号为短号，获取真实房间号: %s ", realId)
	} else {
		realId = roomId
	}

	return &BiliRoom{
		roomId:  roomId,
		realId:  realId,
		conn:    conn,
		recMsg:  recMsg,
		OutMsg:  OutMsg,
		isAlive: false,
		timeout: 3,
		MsgHandler: func(mi *model.MessageInfo) error {
			ulog.Infof("[弹幕] %s: %s", mi.Info.([]interface{})[2].([]interface{})[1], mi.Info.([]interface{})[1])
			return nil
		},
	}
}

func (b *BiliRoom) Start() {
	for {

		if !b.isAlive {
			if err := b.getRoomInfo(); err != nil {
				ulog.Info("房间信息获取失败:", err)
				goto reconnect
			}
			if err := b.connect(); err != nil {
				ulog.Info("房间连接失败:", err)
				goto reconnect
			}
			if err := b.verify(); err != nil {
				ulog.Info("房间验证失败:", err)
				goto reconnect
			}
			go b.readMessage()
			go b.decodeMsg()
			go b.heartBeat()
		}
	reconnect:
		time.Sleep(time.Second * time.Duration(b.timeout))
	}
}

func (b *BiliRoom) connect() error {
	var err error
	b.conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("wss://%s/%s", b.address, "sub"), nil)
	if err != nil {
		return err
	}
	b.isAlive = true
	return nil
}

func (b *BiliRoom) verify() error {
	// 发送房间验证包
	roomInfo := fmt.Sprintf(`{"uid": 0, "roomid": %s, "protover": 3, "platform": "web", "type": 2, "key": "%s"}`, b.realId, b.token)
	err := b.conn.WriteMessage(websocket.BinaryMessage, __pack(roomInfo, 1, 7))
	if err != nil {
		b.isAlive = false
		ulog.Info("write:", err)
		return err
	}
	_, _, err = b.conn.ReadMessage()
	if err != nil {
		b.isAlive = false
		return err
	}
	ulog.Info("房间连接成功")
	return nil
}

func (b *BiliRoom) readMessage() error {
	for {
		if !b.isAlive {
			return errors.New("房间已断开连接！")
		}
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			ulog.Info("read:", err)
			b.isAlive = false
			return err
		}
		b.recMsg <- message
	}
}

func (b *BiliRoom) decodeMsg() {
	for msg := range b.recMsg {
		if !b.isAlive {
			return
		}
		utils.DecodeMessage(msg, b.OutMsg)
	}
}

func (b *BiliRoom) heartBeat() error {
	// 心跳包
	for {
		time.Sleep(time.Second * time.Duration(30))
		err := b.conn.WriteMessage(websocket.BinaryMessage, __pack("[object Object]", 1, 2))
		if err != nil {
			ulog.Info(err)
			b.isAlive = false
			return err
		}
		ulog.Info("[心跳包] 发送成功")

	}
}

func (b *BiliRoom) getRoomInfo() error {
	// 获取房间弹幕地址
	ri := new(model.RoomInfo)
	res, err := http.Get("https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=" + b.realId)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &ri)
	if ri.Code != 0 || len(ri.Data.HostList) == 0 {
		err = errors.New("获取房间信息失败！")
	}
	if err != nil {
		return err
	}
	b.address = ri.Data.HostList[0].Host
	b.token = ri.Data.Token
	return nil
}

func getRealId(roomId string) (string, error) {
	// 真实房间号
	ri := new(model.RealIdInfo)
	res, err := http.Get("https://api.live.bilibili.com/xlive/web-room/v1/index/getRoomPlayInfo?room_id=" + roomId)
	if err != nil {
		return roomId, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	err = json.Unmarshal(data, &ri)
	if ri.Code != 0 {
		err = errors.New("获取房间信息失败！")
	}
	if err != nil {
		return roomId, err
	}
	realId := strconv.Itoa(ri.Data.RoomID)
	return realId, nil
}
func __pack(s string, i int, j int) []byte {
	// 字节流打包
	format := []string{"H", "H", "I", "I"}
	values := []interface{}{16, i, j, 1}
	bp := new(utils.BinaryPack)
	data, _ := bp.Pack(format, values)
	data = append(data, []byte(s)...)
	bp2 := new(utils.BinaryPack)
	data2, _ := bp2.Pack([]string{"I"}, []interface{}{len(data) + 4})
	data = append(data2, data...)
	return data
}

func (dm *BiliRoom) DanmuHandler() {
	ulog.Info("开始处理弹幕命令")
	var danmu = map[int]string{}
	var scList = map[int]model.SuperChatInfo{}
	for message := range dm.OutMsg {
		mi := new(model.MessageInfo)
		json.Unmarshal(message, &mi)
		switch mi.Cmd {
		case "POP":
			// 人气
			pi := new(model.PopInfo)
			json.Unmarshal(message, &pi)
			// ulog.Infof("[人气] %d", pi.Count)

		case "DANMU_MSG":
			// 弹幕
			//ulog.Infof("[弹幕] %s: %s", mi.Info.([]interface{})[2].([]interface{})[1], mi.Info.([]interface{})[1])
			danmu[int(mi.Info.([]interface{})[2].([]interface{})[0].(interface{}).(float64))] = mi.Info.([]interface{})[1].(string)
			err := dm.MsgHandler(mi)
			if err != nil {
				ulog.Error(err)
			}
		case "SUPER_CHAT_MESSAGE":
			// SC
			sci := new(model.SuperChatInfo)
			json.Unmarshal(message, &sci)
			ulog.Infof("[SC] %d元 %s: %s %d", sci.Data.Price, sci.Data.UserInfo.Uname, sci.Data.Message, sci.Data.ID)
			scList[int(sci.Data.ID)] = *sci
		case "SUPER_CHAT_MESSAGE_JPN":

		case "SUPER_CHAT_MESSAGE_DELETE":
			// SC被删除
			ulog.Infof("%s", string(message))
			scd := new(model.SuperChatDelete)
			json.Unmarshal(message, &scd)
			sci := scList[int(scd.Data.Ids[0])]
			ulog.Infof("[SC被删除] %s : %s", sci.Data.UserInfo.Uname, sci.Data.Message)

		case "SEND_GIFT":
			// 礼物
			// ulog.Infof("%s", string(message[16:]))
			gf := new(model.GiftInfo)
			json.Unmarshal(message, &gf)
			// ulog.Infof("[礼物] %s 赠送 %d个 %s %.1f元", gf.Data.Uname, gf.Data.Num, gf.Data.GiftName, float64(gf.Data.Price)/1000*float64(gf.Data.Num))
		case "COMBO_SEND":
			// 连击礼物
			// ulog.Infof(string(message))
			ci := new(model.ComboInfo)
			json.Unmarshal(message, &ci)
			// ulog.Infof("[礼物连击] %s 连续赠送 %d个 %s ", ci.Data.Uname, ci.Data.ComboNum, ci.Data.GiftName)
		case "GUARD_BUY":
			// 大航海
		case "USER_TOAST_MSG":
			// 大航海
			ci := new(model.CrewInfo)
			json.Unmarshal(message, &ci)
			// ulog.Infof("[大航海] %s 开通 %s * %d%s %d元", ci.Data.Username, ci.Data.RoleName, ci.Data.Num, ci.Data.Unit, ci.Data.Price/1000)
		case "ONLINE_RANK_V2":

		case "ONLINE_RANK_TOP3":

		case "INTERACT_WORD":

		case "ENTRY_EFFECT":

		case "ROOM_REAL_TIME_MESSAGE_UPDATE":

		case "ONLINE_RANK_COUNT":

		case "HOT_RANK_CHANGED_V2":

		case "LIVE":
			ulog.Infof("开播了")
			// fmt.Println(string(message[16:]))
		case "PREPARING":
			ulog.Infof("已下播")
			danmu = map[int]string{}
			scList = map[int]model.SuperChatInfo{}
			// fmt.Println(string(message[16:]))
		case "ROOM_CHANGE":

		case "WATCHED_CHANGE":

		case "STOP_LIVE_ROOM_LIST":

		case "HOT_ROOM_NOTIFY":

		case "HOT_RANK_CHANGED":

		case "HOT_RANK_SETTLEMENT":

		case "HOT_RANK_SETTLEMENT_V2":

		case "LIVE_INTERACTIVE_GAME":

		case "VOICE_JOIN_LIST":
			// 连麦请求
			// {"cmd":"VOICE_JOIN_LIST","data":{"cmd":"","room_id":23606554,"category":1,"apply_count":6,"red_point":1,"refresh":1},"room_id":23606554}
		case "VOICE_JOIN_ROOM_COUNT_INFO":
			// 连麦消息
			// {"cmd":"VOICE_JOIN_ROOM_COUNT_INFO","data":{"cmd":"","room_id":23606554,"root_status":1,"room_status":1,"apply_count":5,"notify_count":0,"red_point":1},"room_id":23606554}
		case "ROOM_BLOCK_MSG":
			// 禁言个人
			// {"cmd":"ROOM_BLOCK_MSG","data":{"dmscore":30,"operator":2,"uid":1772442517,"uname":"晓小轩iAvA"},"uid":"1772442517","uname":"晓小轩iAvA"}
			bi := new(model.BlockInfo)
			json.Unmarshal(message, &bi)
			ulog.Infof("[禁言] %s 被禁言", bi.Data.Uname)
			ulog.Infof("上一条弹幕是: %s", danmu[bi.Data.UID])
			// ulog.Infof("%s", string(message[16:]))
		case "ROOM_SILENT_ON":
			// 开启禁言
		case "ROOM_SILENT_OFF":
			// 关闭禁言
		case "WIDGET_BANNER":

		case "COMMON_NOTICE_DANMAKU":

		case "ANCHOR_LOT_START":
			// 天选
			// {"cmd":"ANCHOR_LOT_START","data":{"asset_icon":"https://i0.hdslb.com/bfs/live/627ee2d9e71c682810e7dc4400d5ae2713442c02.png","award_image":"","award_name":"年度大会员","award_num":40,"cur_gift_num":0,"current_time":1653042849,"danmu":"哔哩哔哩 (゜-゜)つロ 干杯~","gift_id":31039,"gift_name":"牛哇牛哇","gift_num":1,"gift_price":100,"goaway_time":180,"goods_id":15,"id":2656528,"is_broadcast":1,"join_type":1,"lot_status":0,"max_time":600,"require_text":"关注主播","require_type":1,"require_value":0,"room_id":25059330,"send_gift_ensure":0,"show_panel":1,"start_dont_popup":0,"status":1,"time":599,"url":"https://live.bilibili.com/p/html/live-lottery/anchor-join.html?is_live_half_webview=1\u0026hybrid_biz=live-lottery-anchor\u0026hybrid_half_ui=1,5,100p,100p,000000,0,30,0,0,1;2,5,100p,100p,000000,0,30,0,0,1;3,5,100p,100p,000000,0,30,0,0,1;4,5,100p,100p,000000,0,30,0,0,1;5,5,100p,100p,000000,0,30,0,0,1;6,5,100p,100p,000000,0,30,0,0,1;7,5,100p,100p,000000,0,30,0,0,1;8,5,100p,100p,000000,0,30,0,0,1","web_url":"https://live.bilibili.com/p/html/live-lottery/anchor-join.html"}}
		case "ANCHOR_LOT_AWARD":
			// 中奖名单

		case "POPULARITY_RED_POCKET_NEW":
			// 新红包

		case "POPULARITY_RED_POCKET_START":
			// 红包开始

		case "POPULARITY_RED_POCKET_WINNER_LIST":
			// 红包中奖名单

		case "NOTICE_MSG":
			// 通知消息
			// if strings.Contains(string(message[16:]), "舰长") || strings.Contains(string(message[16:]), "提督") || strings.Contains(string(message[16:]), "总督") {
			// 	return
			// }
		default:

			ulog.Infof("%s", string(message))
		}
	}
	// var w1 sync.WaitGroup
	// w1.Add(1)
	// w1.Wait()
}
