package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//_ "github.com/marinellirubens/dbwrapper"
)

// teste struct to understand how to send the response
type Teste struct {
	status string
}

// contants for the server (probably will move to a file)
const (
	SERVER_ADDRESS string = "localhost"
	SERVER_PORT    int    = 8080
)

// example of method to handle a request
func getInfo(c *gin.Context) {
	//response := Teste{status: "OK"}
	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

// serve the api
func ServeApi() {
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
