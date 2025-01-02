package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// Variable to contain the default path for the configuration file
const DefaultCfgFilePath = "./config.ini"

type Database struct {
	Id       string `json:"dbid"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Dbname   string `json:"dbname"`
	Dbtype   string `json:"dbtype"`
	Service  string `json:"service"`
	Password string `json:"password"`
}

type ConfigJson struct {
	Server struct {
		Server_port    int    `json:"server_port"`
		Server_address string `json:"server_address"`
		Logger_file    string `json:"logger_file"`
	} `json:"server"`
	Databases []Database `json:"databases"`
}

func (s *ConfigJson) PrintInfo() {
	fmt.Printf("Server info %s:%d file: %s\n",
		s.Server.Server_address,
		s.Server.Server_port,
		s.Server.Logger_file,
	)
	//fmt.Printf("%v\n", s.Databases)
	for _, database := range s.Databases {
		fmt.Printf("%v\n", database)
		fmt.Printf("%v\n", database.Password)
	}
}

func GetJsonConfig(filename string) (ConfigJson, error) {
	config := ConfigJson{}
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
		return config, err
	}
	return config, nil
}

// gets the configuration file
func GetInfoFile(filename ...string) (*ini.File, error) {
	var file_path string
	//example of reading a json file
	//jsonFile, err := json.NewDecoder(DefaultCfgFilePath)

	if len(filename) == 0 {
		file_path = DefaultCfgFilePath
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
