package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// DefaultCfgFilePath ... Variable to contain the default path for the configuration file
const DefaultCfgFilePath = "./config.json"

type Database struct {
	ID       string `json:"dbid"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Dbname   string `json:"dbname"`
	Dbtype   string `json:"dbtype"`
	Service  string `json:"service"`
	Password string `json:"password"`
}

type APIKey struct {
	Key        string   `json:"key"`
	AllowedDbs []string `json:"allowedDbs"`
}

type ServerConfig struct {
	ServerPort    int      `json:"server_port"`
	ServerAddress string   `json:"server_address"`
	LoggerFile    string   `json:"logger_file"`
	LogLevel      string   `json:"loglevel"`
	APIKeys       []APIKey `json:"apikeys"`
}

type ConfigJSON struct {
	Server    ServerConfig `json:"server"`
	Databases []Database   `json:"databases"`
}

func (s *ConfigJSON) PrintInfo() {
	fmt.Printf("Server info %s:%d file: %s\n",
		s.Server.ServerAddress,
		s.Server.ServerPort,
		s.Server.LoggerFile,
	)
	//fmt.Printf("%v\n", s.Databases)
	for _, database := range s.Databases {
		fmt.Printf("%v\n", database)
		fmt.Printf("%v\n", database.Password)
	}
}

func GetJSONConfig(filename string) (ConfigJSON, error) {
	fmt.Println("Reading configuration file", filename)
	config := ConfigJSON{}
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
