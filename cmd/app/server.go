package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/marinellirubens/dbwrapper/database"
	"github.com/marinellirubens/dbwrapper/internal/config"
	logs "github.com/marinellirubens/dbwrapper/internal/logger"
)

var signalChan = make(chan os.Signal, 1)

func ServeApiNative(address string, port int, app *database.App) error {

	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()
	errArr := app.SetupDbConnections()
	if errArr != nil {
		app.Log.Fatal("Error connecting to databases")
	}

	handler, err := SetupRoutes(mux, app)
	if err != nil {
		app.Log.Fatal("Error setting up the routes")
		return err
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))

	signal.Notify(signalChan, os.Interrupt)

	go handlesUserInterrupt(app)
	defer handlesPanic(app)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		app.Log.Fatal(fmt.Sprintf("http server error: %s", err))
		return err
	}
	return nil
}

func closeConnections(app *database.App) {
	app.Log.Info("Closing connections")
	for key, database := range app.DbConns {
		app.Log.Info(fmt.Sprintf("Closing connection %v", key))
		err := database.Close()
		if err != nil {
			app.Log.Fatal(fmt.Sprintf("Error trying to close connection %s  %v", key, err))
		}
	}
}

func handlesPanic(app *database.App) {
	if err := recover(); err != nil { //catch
		app.Log.Error("Received a panic signal, stopping service...")
		closeConnections(app)

		fmt.Fprintf(os.Stderr, "Error receiving : %v\n", err)
		os.Exit(1)
	}
}

func handlesUserInterrupt(app *database.App) {
	<-signalChan
	app.Log.Error("Received an interrupt signal, stopping service...")

	closeConnections(app)

	app.Log.Debug("Exiting process")
	os.Exit(1)
}

func catch() { //catch or finally
	if err := recover(); err != nil { //catch
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
	}
}

func RunServer(cfgPath string) {
	defer catch()

	cfgj, err := config.GetJsonConfig(cfgPath)
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}

	logger, err := logs.CreateLogger(
		cfgj.Server.Logger_file,
		logs.DEBUG,
		[]uint16{logs.STREAM_WRITER, logs.FILE_WRITER},
	)
	logger.Debug(
		fmt.Sprintf(
			"Logger setup using config file %s logging into log file %s",
			cfgPath,
			cfgj.Server.Logger_file,
		),
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
			application.Log.Error("Error setting up databases")
			panic(err)
		}
	}

	_ = ServeApiNative(host, port, application)
}

func SetupAppDbs(dbInfo config.Database, app *database.App) error {
	//fmt.Printf("Database: %v\n", dbInfo)
	switch method := dbInfo.Dbtype; method {
	case database.ORACLE:
		db := database.SetOracleConnection(dbInfo)

		app.DbHandlers[dbInfo.Id] = &db
		return nil
	case database.POSTGRES:
		db := database.SetPostgresConnection(dbInfo)

		app.DbHandlers[dbInfo.Id] = &db
		return nil
	case database.MYSQL:
		db := database.SetMysqlConnection(dbInfo)

		app.DbHandlers[dbInfo.Id] = &db
		return nil
	default:
		return errors.New("Database type not found")
	}
}
