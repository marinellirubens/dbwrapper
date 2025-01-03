package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/marinellirubens/dbwrapper/internal/logger"
	"github.com/marinellirubens/dbwrapper/internal/utils"
)

// message for http response
const METHOD_NOT_ALLOWED = "command not allowed on this endpoint"

// Application object to handle the endpoints and connection with database
type App struct {
	// database connection
	Postgres PostgresHandler
	Oracle   OracleHandler
	Mongo    MongoHandler

	DbHandlers map[string]DbConnection
	DbConns    map[string]*sql.DB

	// logger object for general purposes
	Log *logger.Logger
}

func (app *App) SetupDbConnections() {
	var db *sql.DB
	for _, dbInfo := range app.DbHandlers {
		switch handlerType := dbInfo.GetDbType(); handlerType {
		case ORACLE:
			db = GetOracleConnection(dbInfo)
		case POSTGRES:
			db, _ = GetPostgresConnection(dbInfo)
		default:
			app.Log.Warning("Handler not setup")
			return
		}

		handlerType := dbInfo.GetDbId()
		app.DbConns[handlerType] = db
	}
}

func (app *App) GetDatabasesRequest(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		app.Log.Debug(fmt.Sprintf("Requested ping by %v", utils.ReadUserIP(r)))
		app.Log.Debug(fmt.Sprintf("Headers %s", r.Header))

		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(app.DbHandlers)
		_, err := w.Write(jsonResponse)
		if err != nil {
			app.Log.Error(fmt.Sprintf("Error trying to get server. %v", err))
			panic(err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func (app *App) ProcessGenericRequest(w http.ResponseWriter, r *http.Request) {
	// treat the method
	var err error
	var result []byte

	query := r.Header.Get("query")
	dbId := r.Header.Get("database")
	dbConnection, ok := app.DbConns[dbId]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Database not located")); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
	}

	switch method := r.Method; method {

	case http.MethodGet: // process query and return a json
		result, err = getQueryFromDatabase(query, dbConnection, app)
	case http.MethodDelete: // process deletes on the database
		result, err = processDelete(query, dbConnection, app.Log)
	case http.MethodPatch: // process updates on the database
		result, err = processUpdate(query, dbConnection, app.Log)
	case http.MethodPost: // process inserts on the database
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("Method not implemented yet")); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("Method not allowed")); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
	}

	if err != nil {
		errMessage := fmt.Sprintf("Error processing request %v", err)
		app.Log.Error(errMessage)

		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(errMessage)); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
	}
}

// processes the update on postgresql database
func processUpdate(command string, db *sql.DB, Log *logger.Logger) ([]byte, error) {
	start := time.Now()
	err := validateUpdate(command)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	Log.Info(fmt.Sprintf("Command sent: `%s` processing...", command))
	result, err := db.Exec(command)
	Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	if err != nil {
		Log.Error(fmt.Sprintf("%s", err))
		return []byte(fmt.Sprintf("%v\n", err)), err
	} else {
		lastInsert, _ := result.LastInsertId()
		rowsAfected, _ := result.RowsAffected()
		Log.Debug(
			fmt.Sprintf(
				"Process result LastInsertId:%v RowsAffected: %v",
				lastInsert, rowsAfected,
			),
		)
	}

	return []byte("Success"), nil
}

func processDelete(command string, db *sql.DB, Log *logger.Logger) ([]byte, error) {
	start := time.Now()
	err := validateDelete(command)
	if err != nil {
		return nil, err
	}
	Log.Info(fmt.Sprintf("Command sent: `%s` processing...", command))
	_, err = db.Exec(command)
	Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	if err != nil {
		return nil, err
	}
	return []byte("Success"), nil
}

// Requests information from the postgresql database that is connected
//
//	Validates if the method is GET, if the method is not GET, returns a StatusMethodNotAllowed response
func (app *App) getQueryFromPostgres(query string) ([]byte, error) {
	var err error

	start := time.Now()
	err = validateQuery(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}

	err = app.Postgres.checkDbConnection()
	if err != nil {
		return nil, err
	}

	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	rows, err := app.Postgres.db.Query(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	allgeneric := make([]map[string]interface{}, 0)
	colvals := make([]interface{}, len(cols))
	for rows.Next() {
		colassoc := make(map[string]interface{}, len(cols))
		for i := range colvals {
			colvals[i] = new(interface{})
		}
		if err := rows.Scan(colvals...); err != nil {
			return nil, err
		}

		for i, col := range cols {
			colassoc[col] = *colvals[i].(*interface{})
		}
		allgeneric = append(allgeneric, colassoc)
	}

	err2 := rows.Close()
	if err2 != nil {
		return nil, err2
	}

	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	js, _ := json.Marshal(allgeneric)

	return js, nil
}

// Requests information from the database
//
//	Validates if the method is GET, if the method is not GET, returns a StatusMethodNotAllowed response
func getQueryFromDatabase(query string, db *sql.DB, app *App) ([]byte, error) {
	var err error

	start := time.Now()
	app.Log.Debug("Validating query")
	err = validateQuery(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}

	app.Log.Debug("Checking db connection")
	err = checkDbConnection(db)
	if err != nil {
		return nil, err
	}

	app.Log.Debug("Processing query")
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	rows, err := db.Query(query)
	if err != nil {
		app.Log.Error(fmt.Sprintf("Error processing query %v", err))
		return nil, err
	}
	app.Log.Debug("Processing rows")
	cols, err := rows.Columns()
	if err != nil {
		app.Log.Error(fmt.Sprintf("Error processing rows %v", err))
		return nil, err
	}

	allgeneric := make([]map[string]interface{}, 0)
	colvals := make([]interface{}, len(cols))
	for rows.Next() {
		colassoc := make(map[string]interface{}, len(cols))
		for i := range colvals {
			colvals[i] = new(interface{})
		}
		if err := rows.Scan(colvals...); err != nil {
			app.Log.Error(fmt.Sprintf("Error processing rows %v", err))
			return nil, err
		}

		for i, col := range cols {
			colassoc[col] = *colvals[i].(*interface{})
		}
		allgeneric = append(allgeneric, colassoc)
	}

	err2 := rows.Close()
	if err2 != nil {
		app.Log.Error(fmt.Sprintf("Error processing rows %v", err))
		return nil, err2
	}

	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	js, _ := json.Marshal(allgeneric)

	return js, nil
}
