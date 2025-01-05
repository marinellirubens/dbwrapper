package main

import (
	"log"

	cfg "github.com/marinellirubens/dbwrapper/internal/config"
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
	//dbl.SetMapping()
	configJson, err := cfg.GetJsonConfig("./test.json")
	if err != nil {
		log.Fatal(err)
	}
	configJson.PrintInfo()

	//config := Config{}
	//configFile, err := os.Open("./test.json")
	//if err != nil {
	//log.Fatal(err)
	//}
	//defer configFile.Close()
	//jsonParser := json.NewDecoder(configFile)
	//err = jsonParser.Decode(&config)
	//if err != nil {
	//fmt.Printf("error parsing json %v", err)
	//}
	//fmt.Println("port", config.Database.Port)
	//dbl.TestConnection()
}
