package config

import (
	"testing"

	cf "github.com/marinellirubens/dbwrapper/internal/config"
)

func Test_GetInfoFile(t *testing.T) {
	//fmt.Println(os.Getwd())
	_, err := cf.GetInfoFile("./examples/config.example.ini")
	if err != nil {
		t.Errorf("GetInfoFile - ERROR %v", err)
	} else {
		t.Logf("GetInfoFile - SUCCESS")
	}
}
