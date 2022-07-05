/*
 * @Author: ww
 * @Date: 2022-07-02 02:25:38
 * @Description:
 * @FilePath: /danmuplay/danmu/danmu_test.go
 */
package danmu

import (
	"testing"
)


func TestNewBiliRoom(t *testing.T) {
	
	tmp := NewBiliRoom("746504")
	go tmp.Start()
	tmp.DanmuHandler()
}

