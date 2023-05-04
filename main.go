package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	//_ "github.com/marinellirubens/dbwrapper"
)

// teste struct to understand how to send the response
// Atributes that need to appear as key in a json needs to be start with capital letters
type Teste struct {
	Information string `json:"information" binding:"required"`
	Teste       string `json:"teste" binding:"required"`
}

// example of method to handle a request
func (app *app) getInfo(c *gin.Context) {
	response := Teste{Information: "OK", Teste: "klsdhkdhfgh"}

	c.Bind(&response)
	fmt.Println(response)

	c.JSON(http.StatusOK, response)
}

// serve the api
func ServeApi(address string, port int, app *app) {
	gin.SetMode(gin.DebugMode)
	server_path := fmt.Sprintf("%v:%v", address, port)
	fmt.Printf("Starting server on %v\n", server_path)

	// define the endpoints/handlers of the api
	router := gin.Default()
	router.GET("/", app.getInfo)
	//router.POST("/albums", postAlbums)
	router.Run(server_path)
}

func main() {
	cfg, err := GetInfoFile("./config.ini")
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}

	fmt.Println(cfg.Section("SERVER"))
	host := cfg.Section("SERVER").Key("SERVER_ADDRESS").String()
	port, _ := cfg.Section("SERVER").Key("SERVER_PORT").Int()

	psqlInfom := GetConnectionInfo(cfg)
	fmt.Println(psqlInfom)

	db, err := ConnectToPsql(psqlInfom)
	if err != nil {
		panic(err)
	}
	application := &app{db: db}

	ServeApi(host, port, application)
}
