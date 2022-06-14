/*
 * @Author: ww
 * @Date: 2022-06-15 02:26:16
 * @Description:
 * @FilePath: /live-informed/bilibili/api_test.go
 */
package bilibili

import (
	"testing"
)

func Test_getLiveInfo(t *testing.T) {
	type args struct {
		uids []int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"case1", args{[]int64{672328094,1265680561}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLiveInfo(tt.args.uids)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLiveInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}
