/*
 * @Author: ww
 * @Date: 2022-06-17 02:04:35
 * @Description:
 * @FilePath: /live-informed/bilibili/entity.go
 */
package bilibili

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

type UserInfoResp struct {
	Code uint64 `json:"code"` // 0：成功 -400：请求错误
	Msg	string `json:"message"`
	Ttl int `json:"ttl"`
	Data UserInfo `json:"data"`
}

type UserInfo struct {
	Mid int `json:"mid"`
	Name string `json:"name"`
	Sex string `json:"sex"`
	Face string `json:"face"`
	FaceNft int `json:"face_nft"`
	FaceNftType int `json:"face_nft_type"`
	Sign string `json:"sign"`
	Rank int `json:"rank"`
	Level int `json:"level"`
	Jointime int `json:"jointime"`
	Moral int `json:"moral"`
	Silence int `json:"silence"`
	Coins int `json:"coins"`
	FansBadge bool `json:"fans_badge"`
	//FansMedal FansMedal `json:"fans_medal"`
	//Official Official `json:"official"`
	//Vip Vip `json:"vip"`
	//Pendant Pendant `json:"pendant"`
	//Nameplate Nameplate `json:"nameplate"`
	//UserHonourInfo UserHonourInfo `json:"user_honour_info"`
	IsFollowed bool `json:"is_followed"`
	TopPhoto string `json:"top_photo"`
	//Theme Theme `json:"theme"`
	//SysNotice SysNotice `json:"sys_notice"`
	LiveRoom *LiveRoom `json:"live_room"`
	Birthday string `json:"birthday"`
	//School School `json:"school"`
	//Profession Profession `json:"profession"`
	Tags interface{} `json:"tags"`
	//Series Series `json:"series"`
	IsSeniorMember int `json:"is_senior_member"`
}

type LiveRoom struct {
	RoomStatus int `json:"roomStatus"`
	LiveStatus int `json:"liveStatus"` // 0：未开播 1：直播中
	URL string `json:"url"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Roomid int `json:"roomid"`
	RoundStatus int `json:"roundStatus"`
	BroadcastType int `json:"broadcast_type"`
	//WatchedShow WatchedShow `json:"watched_show"`
}