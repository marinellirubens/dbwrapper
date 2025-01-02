package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/marinellirubens/dbwrapper/database"
	"github.com/marinellirubens/dbwrapper/internal/config"
	logs "github.com/marinellirubens/dbwrapper/internal/logger"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func ServeApiNative(address string, port int, app *database.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	handler, err := SetupRoutes(mux, app)
	if err != nil {
		panic("Error setting up the routes")
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func RunServer(cfgPath string) {
	cfg, err := config.GetInfoFile(cfgPath)
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}
	cfgj, err := config.GetJsonConfig("./test.json")
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}
	logger, err := logs.CreateLogger(
		cfgj.Server.Logger_file,
		logs.DEBUG,
		[]uint16{logs.STREAM_WRITER, logs.FILE_WRITER},
	)
	if err != nil {
		log.Fatal(err)
	}

	host := cfgj.Server.Server_address
	port := cfgj.Server.Server_port

	psqlInfom := database.GetConnectionInfo(cfg)
	db, err := database.ConnectToPsql(psqlInfom)
	if err != nil {
		panic(err)
	}
	defer database.CloseConn(db)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	application := &database.App{Log: logger, DbHandlers: map[string]database.DbConnection{}}

	for _, v := range cfgj.Databases {
		SetupAppDbs(v, application)
	}

	application.IncludeDbConnection(db, reflect.TypeOf(database.PostgresHandler{}), psqlInfom)
	ServeApiNative(host, port, application)
}

func SetupAppDbs(dbInfo config.Database, app *database.App) error {
	fmt.Printf("Database: %v\n", dbInfo)
	switch method := dbInfo.Dbtype; method {
	case database.ORACLE:
		db := database.SetOracleConnection(dbInfo)

		app.DbHandlers[dbInfo.Id] = &db
		return nil
	case database.POSTGRES:
		db := database.SetPostgresConnection(dbInfo)

		app.DbHandlers[dbInfo.Id] = &db
		return nil
	default:
		return errors.New("Database type not found")
	}
}
