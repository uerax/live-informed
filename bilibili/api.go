/*
 * @Author: ww
 * @Date: 2022-06-15 02:26:16
 * @Description:
 * @FilePath: \live-informed\bilibili\api.go
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
	// 疑似被b站关闭
	StatusInfoUrl = "http://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids"
	UserInfoUrl   = "https://api.bilibili.com/x/space/acc/info?mid="
)

func GetLiveInfo(uids []int64) (*LiveInfoResp, error) {

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

func GetUserInfo(mid string) (*UserInfo, error) {
	userInfo, err := http.Get(UserInfoUrl + mid)
	if err != nil {
		log.Errorf("bilibili接口异常 : %v", err)
		return nil, err
	}

	defer userInfo.Body.Close()
	body, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		log.Errorf("bilibili数据异常 : %v", err)
		return nil, err
	}

	var tmp UserInfoResp
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		log.Errorf("bilibili数据异常 : %v", err)
		return nil, err
	}

	return &tmp.Data, nil
}

func UserIsLiving(mid string) (bool, error) {
	userInfo, err := GetUserInfo(mid)
	if err != nil {
		return false, err
	}
	if userInfo == nil || userInfo.LiveRoom == nil {
		return false, nil
	}
	return userInfo.LiveRoom.LiveStatus == 1, nil
}
