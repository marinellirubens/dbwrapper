package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// TODO: include the possibility to get the configuration file from a different path

// gets the configuration file
func GetInfoFile(filename ...string) (*ini.File, error) {
	var file_path string
	if len(filename) == 0 {
		file_path = "./config.ini"
	} else {
		file_path = filename[0]
	}

	cfg, err := ini.Load(file_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return nil, err
	}
	return cfg, nil
}
