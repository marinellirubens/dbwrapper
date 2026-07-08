package config

import (
	"testing"
)

func Test_GetJSONConfig(t *testing.T) {
	//fmt.Println(os.Getwd())
	_, err := GetJSONConfig("./examples/config.example.json")
	if err != nil {
		t.Errorf("GetInfoFile - ERROR %v", err)
	} else {
		t.Logf("GetInfoFile - SUCCESS")
	}
}
