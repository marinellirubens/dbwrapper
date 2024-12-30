package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		User     string `json:"user"`
		Port     int    `json:"port"`
	} `json:"database"`
	Token string `json:"token"`
}

func main() {
	config := Config{}

	configFile, err := os.Open("./test.json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	fmt.Println(config.Database.Port)
}
