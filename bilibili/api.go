/*
 * @Author: ww
 * @Date: 2022-06-15 02:26:16
 * @Description:
 * @FilePath: /live-informed/bilibili/api.go
 */
package bilibili

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tencent-connect/botgo/log"
)

const (
	StatusInfoUrl = "http://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids"
)

type LiveInfoReq struct {
	Uids []int64 `json:"uids"`
}

type LiveInfoResp struct {
	Code int `json:"code"` // 0：成功 -111：csrf校验失败
	Msg	string `json:"msg"`
	Message string `json:"message"`
	Data map[string]LiveInfo `json:"data"`
}

type LiveInfo struct {
	Title string `json:"title"`	// 直播间标题	
	RoomId int64 `json:"room_id"` // 直播间实际房间号
	Uid int64 `json:"uid"`	// 主播mid	
	Online int64 `json:"online"`	// 直播间在线人数	
	LiveTime int64 `json:"live_time"`	// 直播持续时长	
	LiveStatus int `json:"live_status"` // 0：未开播 1：正在直播 2：轮播中
	ShortId int64 `json:"short_id"` // 直播间短id
	Area int64 `json:"area"`	// 直播间分区id	
	AreaName string `json:"area_name"`	// 直播间分区名	
	AreaV2Id int64 `json:"area_v2_id"`	// 直播间新版分区id	
	Areav2Name string `json:"area_v2_name"`	// 直播间新版分区名	
	AreaV2ParentId int64 `json:"area_v2_parent_id"` // 直播间父分区id	
	Areav2ParentName string `json:"area_v2_parent_name"` // 直播间父分区名	
	UName string `json:"uname"`	// 主播用户名	
	Face string `json:"face"`	// 主播头像url	
	TagName string `json:"tag_name"`	// 直播间标签	
	Tags string `json:"tags"`	// 直播间自定标签	
	CoverFromUser string `json:"cover_from_user"`	// 直播间封面url	
	Keyframe string `json:"keyframe"` // 直播间关键帧url	
	LockTill string `json:"lock_till"`	// 直播间封禁信息	
	HiddenTill string `json:"hidden_till"`	// 直播间隐藏信息	
	BroadcastType int `json:"broadcast_type"`	// 0:普通直播	1：手机直播
}

func getLiveInfo(uids []int64) (*LiveInfoResp, error) {
	
	liveInfoReq := &LiveInfoReq{uids}

	req, err := json.Marshal(liveInfoReq)
	if err != nil {
		log.Errorf("参数数据异常 : %v", err)
		return nil, err
	}

	res, err := http.Post(StatusInfoUrl, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Errorf("bilibili接口异常 : %v", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("bilibili数据异常 : %v", err)
		return nil, err
	}

	var tmp *LiveInfoResp

	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Errorf("bilibili数据异常 : %v", err)
		return nil, err
	}

	return tmp, nil
}