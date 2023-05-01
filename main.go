package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	//_ "github.com/marinellirubens/dbwrapper"
)

type Teste struct {
	status string
}

const (
	SERVER_ADDRESS string = "localhost"
	SERVER_PORT    int    = 8080
)

func getInfo(c *gin.Context) {
	//response := Teste{status: "OK"}
	c.JSON(http.StatusOK, gin.H{"data": "hello world"})
}

func ServeApi() {
	server_path := fmt.Sprintf("%v:%v", SERVER_ADDRESS, SERVER_PORT)
	fmt.Printf("Starting server on %v\n", server_path)

	router := gin.Default()
	router.GET("/", getInfo)
	//router.POST("/albums", postAlbums)
	router.Run(server_path)
}

func main() {
	ServeApi()
}
