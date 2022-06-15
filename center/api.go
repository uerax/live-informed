/*
 * @Author: ww
 * @Date: 2022-06-15 08:27:32
 * @Description:
 * @FilePath: /live-informed/center/api.go
 */
package center

import "strconv"

func AddTask(uid string) error {
	id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return err
	}
	tasks.AddTask(id)
	return nil
}

func TaskStart() {
	tasks.Detection()
}