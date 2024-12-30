package config

import (
	"testing"
)

func Test_GetInfoFile(t *testing.T) {
	//fmt.Println(os.Getwd())
	_, err := GetInfoFile("./examples/config.example.ini")
	if err != nil {
		t.Errorf("GetInfoFile - ERROR %v", err)
	} else {
		t.Logf("GetInfoFile - SUCCESS")
	}
}
