package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//_ "github.com/marinellirubens/dbwrapper"
)

// teste struct to understand how to send the response
// Atributes that need to appear as key in a json needs to be start with capital letters
type Teste struct {
	Information string `json:"information" binding:"required"`
	Teste       string `json:"teste" binding:"required"`
}

// contants for the server (probably will move to a file)
const (
	SERVER_ADDRESS string = "localhost"
	SERVER_PORT    int    = 8080
)

// example of method to handle a request
func getInfo(c *gin.Context) {
	response := Teste{Information: "OK", Teste: "klsdhkdhfgh"}

	c.Bind(&response)
	fmt.Println(response)

	c.JSON(http.StatusOK, response)
}

// serve the api
func ServeApi() {
	gin.SetMode(gin.DebugMode)
	server_path := fmt.Sprintf("%v:%v", SERVER_ADDRESS, SERVER_PORT)
	fmt.Printf("Starting server on %v\n", server_path)

	// define the endpoints/handlers of the api
	router := gin.Default()
	router.GET("/", getInfo)
	//router.POST("/albums", postAlbums)
	router.Run(server_path)
}

func main() {
	ServeApi()
}
