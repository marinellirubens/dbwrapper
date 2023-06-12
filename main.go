package main

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
	"reflect"

	cf "github.com/marinellirubens/dbwrapper/config"
	pg "github.com/marinellirubens/dbwrapper/database"
	logs "github.com/marinellirubens/dbwrapper/logger"
)

// TODO: need to create some treatment on the path variable to understand how to do that without any framework

func ServeApiNative(address string, port int, app *pg.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { fmt.Println("Root handler:", r.URL.Path) }))

	mux.HandleFunc("/pg", app.ProcessPostgresRequest)                                                    // process base request for postgresl
	mux.Handle("/pg/", http.StripPrefix("/pg/", http.HandlerFunc(app.ProcessPostgresRequestHandlePath))) // process requests with path arguments

	mux.HandleFunc("/oracle", app.ProcessOracleRequest)
	mux.HandleFunc("/mongodb", app.ProcessMongoRequest)

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
	application := &pg.App{Log: logger}
	application.IncludeDbConnection(db, reflect.TypeOf(pg.PostgresHandler{db: nil, connection_string: psqlInfom}))
	ServeApiNative(host, port, application)
}
