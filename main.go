package main

// TODO: Include the possibility to use a different path for the ini file
// TODO: include connection with redis
// TODO: include connection with mongodb
// TODO: include connection with Mysql
// TODO: Include cli arguments using cobra or urfave/cli
// TODO: include authentication using JWT
// TODO: implement cli arguments
// TODO: improve the readme with examples
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

	mux.Handle("/", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		fmt.Println("Root handler:", r.URL.Path)
	}))

	// process base request for postgresl
	mux.HandleFunc(
		"/pg",
		app.ProcessPostgresRequest,
	)

	// process requests with path arguments
	mux.Handle(
		"/pg/",
		http.StripPrefix("/pg/", http.HandlerFunc(app.ProcessPostgresRequestHandlePath)),
	)

	mux.HandleFunc("/oracle", app.ProcessOracleRequest)
	mux.HandleFunc("/mongodb", app.ProcessMongoRequest)

	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))
	http.ListenAndServe(server_path, mux)
}

func main() {
	logger, err := logs.CreateLogger("/tmp/server.log", logs.DEBUG, []uint16{logs.STREAM_WRITER, logs.FILE_WRITER})
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
	if err := db.Ping(); err != nil {
		panic(err)
	}
	application := &pg.App{Log: logger}
	application.IncludeDbConnection(db, reflect.TypeOf(pg.PostgresHandler{}), psqlInfom)
	ServeApiNative(host, port, application)
}
