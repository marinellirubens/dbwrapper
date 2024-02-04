package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	logs "github.com/marinellirubens/dbwrapper/logger"
)

const METHOD_NOT_ALLOWED = "command not allowed on this endpoint"

// Application object to handle the endpoints and connection with database
type App struct {
	// database connection
	Postgres PostgresHandler
	Oracle   OracleHandler
	Mongo    MongoHandler
	// logger object for general purposes
	Log *logs.Logger
}

func (app *App) IncludeDbConnection(db *sql.DB, handler reflect.Type, connection_string string) {
	app.Log.Info(handler.String())

	switch handlerType := handler.String(); handlerType {
	case "database.PostgresHandler":
		app.Log.Debug("Including postgres handler")
		app.Postgres = PostgresHandler{db: db, connection_string: connection_string}
	case "database.OracleHandler":
		app.Log.Debug("Including oracle handler")
		app.Oracle = OracleHandler{db: db, connection_string: connection_string}
	case "database.MongoHandler":
		app.Log.Debug("Including mongo handler")
		app.Mongo = MongoHandler{db: db, connection_string: connection_string}
	default:
		app.Log.Warning("Handler not setup")
	}
}

// TODO: treat the path to use the second element as the name of the database
func (app *App) ProcessOracleRequest(w http.ResponseWriter, r *http.Request) {
	if app.Oracle.db == nil {
		app.Log.Warning("No Oracle handler was setup")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(nil))
		return
	}
}

// TODO: treat the path to use the second element as the name of the database
func (app *App) ProcessMongoRequest(w http.ResponseWriter, r *http.Request) {
	if app.Mongo.db == nil {
		app.Log.Warning("No Mongodb handler was setup")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(nil))
		return
	}
}

// process selects (GET), delete(DELETE) and update(PATCH)
func (app *App) ProcessPostgresRequestHandlePath(w http.ResponseWriter, r *http.Request) {
	// method to handle path variables
	fmt.Println("Path handler:", r.URL.Path)
	fmt.Println("Path handler:", r.URL.Query().Get("query"))
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(METHOD_NOT_ALLOWED))
}

// process selects (GET), delete(DELETE) and update(PATCH)
func (app *App) ProcessPostgresRequest(w http.ResponseWriter, r *http.Request) {
	// method to handle path variables
	if app.Postgres.db == nil {
		app.Log.Warning("No posgresql handler was setup")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(nil))
		return
	}

	var status int
	var result []byte
	var err error

	// treat the method
	switch method := r.Method; method {
	case http.MethodGet:
		query := r.URL.Query().Get("query")
		result, err = app.getQueryFromPostgres(query)
	case http.MethodDelete:
		query := r.URL.Query().Get("query")
		result, err = app.deleteFromPostgres(query)
	case http.MethodPatch:
		query := r.URL.Query().Get("query")
		result, err = app.updatePostgres(query)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	if err != nil {
		status = http.StatusBadRequest
	} else {
		status = http.StatusOK
	}

	w.WriteHeader(status)
	w.Write(result)
}

// processes the update on postgresql database
func (app *App) updatePostgres(query string) ([]byte, error) {
	start := time.Now()
	err := validateUpdate(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	_, err = app.Postgres.db.Exec(query)
	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	return []byte("Success"), nil
}

// process the delete process on postgresql database
func (app *App) deleteFromPostgres(query string) ([]byte, error) {
	start := time.Now()
	err := validateDelete(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	_, err = app.Postgres.db.Exec(query)
	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
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
		panic(err)
	}

	allgeneric := make([]map[string]interface{}, 0)
	colvals := make([]interface{}, len(cols))
	for rows.Next() {
		colassoc := make(map[string]interface{}, len(cols))
		for i := range colvals {
			colvals[i] = new(interface{})
		}
		if err := rows.Scan(colvals...); err != nil {
			panic(err)
		}

		for i, col := range cols {
			colassoc[col] = *colvals[i].(*interface{})
		}
		allgeneric = append(allgeneric, colassoc)
	}

	err2 := rows.Close()
	if err2 != nil {
		panic(err2)
	}

	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	js, _ := json.Marshal(allgeneric)

	return js, nil
}
