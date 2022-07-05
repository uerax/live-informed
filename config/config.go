/*
 * @Author: ww
 * @Date: 2022-06-15 01:29:40
 * @Description:
 * @FilePath: /live-informed/config/config.go
 */

package config

import (
	"io/ioutil"
	"strconv"

	"github.com/tencent-connect/botgo/log"
	"gopkg.in/yaml.v3"
)

var (
	url = "./etc/live-informed.yaml"
	cfg = make(map[string]map[string]string, 0)
)

func Init() {
	loadConfig()
}

func Reload() {
	loadConfig()
}

func loadConfig() error {
	data, err := ioutil.ReadFile(url)
	if err != nil {
		log.Errorf("配置文件读取失败 : %v", err)
		return err
	}

	var tmp map[string]map[string]string

	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		log.Errorf("配置文件解析错误 : %v", err)
		return err
	}

	cfg = tmp

	return nil

}

func GetToken() string {
	return cfg["server"]["token"]
}

func GetAppId() uint64 {
	appid, _ := strconv.ParseUint(cfg["server"]["appid"], 10, 64)
	return appid
}

func GetRoomId() string {
	return cfg["server"]["roomid"]
}
