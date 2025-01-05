package config

import (
	"testing"
)

func Test_GetJsonConfig(t *testing.T) {
	//fmt.Println(os.Getwd())
	_, err := GetJsonConfig("./examples/config.example.json")
	if err != nil {
		t.Errorf("GetInfoFile - ERROR %v", err)
	} else {
		t.Logf("GetInfoFile - SUCCESS")
	}
}
