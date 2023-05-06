package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cf "github.com/marinellirubens/dbwrapper/config"
	pg "github.com/marinellirubens/dbwrapper/postgres"

	"github.com/gin-gonic/gin"
	//_ "github.com/marinellirubens/dbwrapper"
)

// teste struct to understand how to send the response
// Atributes that need to appear as key in a json needs to be start with capital letters

// serve the api
func ServeApi(address string, port int, app *pg.App) {
	gin.SetMode(gin.ReleaseMode)
	server_path := fmt.Sprintf("%v:%v", address, port)
	fmt.Printf("Starting server on %v\n", server_path)

	// define the endpoints/handlers of the api
	router := gin.Default()
	router.GET("/", app.GetInfo)
	router.GET("/pg", app.GetInfoFromDb)
	//router.POST("/albums", postAlbums)
	router.Run(server_path)
}

func ServeApiNative(address string, port int, app *pg.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.GetInfoNative)
	mux.HandleFunc("/teste", app.GetInfoNativeTeste)

	fmt.Printf("Starting server on %v\n", server_path)
	http.ListenAndServe(server_path, mux)
}

func main() {
	cfg, err := cf.GetInfoFile("./config/config.ini")
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}

	//fmt.Println(cfg.Section("SERVER"))
	host := cfg.Section("SERVER").Key("SERVER_ADDRESS").String()
	port, _ := cfg.Section("SERVER").Key("SERVER_PORT").Int()

	psqlInfom := pg.GetConnectionInfo(cfg)
	//fmt.Println(psqlInfom)

	db, err := pg.ConnectToPsql(psqlInfom)
	if err != nil {
		panic(err)
	}
	application := &pg.App{Db: db}

	ServeApiNative(host, port, application)
}
