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

	"github.com/marinellirubens/dbwrapper/src/database"
	"github.com/marinellirubens/dbwrapper/src/internal/config"
	logs "github.com/marinellirubens/dbwrapper/src/internal/logger"
)

var signalChan = make(chan os.Signal, 1)

func ServeAPINative(address string, port int, app *database.App) error {

	serverPath := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()
	errArr := app.SetupDBConnections()
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
	app.Log.Info(fmt.Sprintf("Starting server on %v", serverPath))

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

	cfgj, err := config.GetJSONConfig(cfgPath)
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}

	logger, err := logs.CreateLogger(
		cfgj.Server.LoggerFile,
		cfgj.Server.LogLevel,
		[]uint16{logs.StreamWriter, logs.FileWriter},
	)
	logger.Debug(
		fmt.Sprintf(
			"Logger setup using config file %s logging into log file %s",
			cfgPath,
			cfgj.Server.LoggerFile,
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	host := cfgj.Server.ServerAddress
	port := cfgj.Server.ServerPort
	application := &database.App{
		Log:        logger,
		DBHandlers: map[string]database.DbConnection{},
		DbConns:    make(map[string]*sql.DB),
		Config:     &cfgj.Server,
	}

	for _, v := range cfgj.Databases {
		err = SetupAppDbs(v, application)
		if err != nil {
			application.Log.Error("Error setting up databases")
			panic(err)
		}
	}

	_ = ServeAPINative(host, port, application)
}

func SetupAppDbs(dbInfo config.Database, app *database.App) error {
	//fmt.Printf("Database: %v\n", dbInfo)
	switch method := dbInfo.Dbtype; method {
	case database.ORACLE:
		db := database.SetOracleConnection(dbInfo)

		app.DBHandlers[dbInfo.ID] = &db
		return nil
	case database.POSTGRES:
		db := database.SetPostgresConnection(dbInfo)

		app.DBHandlers[dbInfo.ID] = &db
		return nil
	case database.MYSQL:
		db := database.SetMysqlConnection(dbInfo)

		app.DBHandlers[dbInfo.ID] = &db
		return nil
	default:
		return errors.New("database type not found")
	}
}
