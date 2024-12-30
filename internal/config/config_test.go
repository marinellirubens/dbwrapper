package config

import (
	"testing"
)

func Test_GetInfoFile(t *testing.T) {
	_, err := GetInfoFile("./config.ini")
	if err != nil {
		t.Errorf("GetInfoFile - ERROR %v", err)
	} else {
		t.Logf("GetInfoFile - SUCCESS")
	}
}
