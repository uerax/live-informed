/*
 * @Author: ww
 * @Date: 2022-06-15 01:53:07
 * @Description:
 * @FilePath: /live-informed/config/config_test.go
 */
package config

import "testing"


func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"case1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("%v", cfg)
		})
	}
}
