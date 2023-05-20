package main

//TODO: include endpoint to insert and delete from database
//TODO: include connection with redis
//TODO: include connection with mongodb
//TODO: include connection with Oracle
//TODO: include connection with Mysql
//TODO: include authentication using JWT
//TODO: implement cli arguments
//TODO: improve the readme with examples
import (
	"fmt"
	"log"
	"net/http"
	"os"

	cf "github.com/marinellirubens/dbwrapper/config"
	logs "github.com/marinellirubens/dbwrapper/logger"
	pg "github.com/marinellirubens/dbwrapper/postgres"
)

func ServeApiNative(address string, port int, app *pg.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	mux.HandleFunc("/pg", app.ProcessPostgresRequest)

	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))
	http.ListenAndServe(server_path, mux)
}

func main() {
	logger, err := logs.CreateLogger("server.log", logs.DEBUG, logs.STREAM_WRITER)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := cf.GetInfoFile("./config/config.ini")
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}

	host := cfg.Section("SERVER").Key("SERVER_ADDRESS").String()
	port, _ := cfg.Section("SERVER").Key("SERVER_PORT").Int()

	psqlInfom := pg.GetConnectionInfo(cfg)
	db, err := pg.ConnectToPsql(psqlInfom)
	if err != nil {
		panic(err)
	}
	application := &pg.App{Db: db, Log: logger}

	ServeApiNative(host, port, application)
}
