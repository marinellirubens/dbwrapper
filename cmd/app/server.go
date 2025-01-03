package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/marinellirubens/dbwrapper/database"
	"github.com/marinellirubens/dbwrapper/internal/config"
	logs "github.com/marinellirubens/dbwrapper/internal/logger"
)

func ServeApiNative(address string, port int, app *database.App) error {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()
	app.SetupDbConnections()

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
	return nil
}

func catch() { //catch or finally
	if err := recover(); err != nil { //catch
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

func RunServer(cfgPath string) {
	//cfg, err := config.GetInfoFile(cfgPath)
	//if err != nil {
	//log.Fatal("error processing configuration file")
	//os.Exit(1)
	//}
	defer catch()

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
	application := &database.App{Log: logger, DbHandlers: map[string]database.DbConnection{}, DbConns: make(map[string]*sql.DB)}

	for _, v := range cfgj.Databases {
		err = SetupAppDbs(v, application)
		if err != nil {
			panic(err)
		}
	}

	//application.IncludeDbConnection(db, reflect.TypeOf(database.PostgresHandler{}), psqlInfom)
	_ = ServeApiNative(host, port, application)
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
