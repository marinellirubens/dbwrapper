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
	DbHandlers map[string]DbConnection
	DbConns    map[string]*sql.DB

	// logger object for general purposes
	Log *logger.Logger
}

func (app *App) SetupDbConnections() []error {
	var db *sql.DB
	var err error
	var errArr []error
	var setupHandle bool = true

	for _, dbInfo := range app.DbHandlers {
		setupHandle = true

		switch handlerType := dbInfo.GetDbType(); handlerType {
		case ORACLE:
			db, err = GetConnection(dbInfo, app.Log)
			if err != nil {
				app.Log.Error(fmt.Sprintf("Error connecting to db %v\n", err))
				errArr = append(errArr, err)
				setupHandle = false
			}
		case POSTGRES:
			db, err = GetConnection(dbInfo, app.Log)
			if err != nil {
				app.Log.Error(fmt.Sprintf("Error connecting to db %v\n", err))
				errArr = append(errArr, err)
				setupHandle = false
			}
		case MYSQL:
			db, err = GetConnection(dbInfo, app.Log)
			if err != nil {
				app.Log.Error(fmt.Sprintf("Error connecting to db %v\n", err))
				errArr = append(errArr, err)
				setupHandle = false
			}
		default:
			app.Log.Warning(fmt.Sprintf("Handler not setup %s", handlerType))
			errArr = append(errArr, fmt.Errorf("Handler not setup"))
			setupHandle = false
		}

		handlerType := dbInfo.GetDbId()
		if setupHandle == true {
			app.DbConns[handlerType] = db
			app.Log.Info(fmt.Sprintf("Setando handler %s", handlerType))
		} else {
			delete(app.DbHandlers, handlerType)
			app.Log.Warning(fmt.Sprintf("Deletando handler %s", handlerType))
		}
	}
	handlers, _ := json.Marshal(app.DbHandlers)
	app.Log.Warning(fmt.Sprintf("Database Handlers %s", handlers))
	return errArr
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
			app.Log.Error(fmt.Sprintf("Error trying to write response. %v", err))
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := w.Write([]byte("Method not allowed"))
		if err != nil {
			app.Log.Error(fmt.Sprintf("Error trying to write response. %v", err))
		}
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
		if _, err := w.Write([]byte("Database not located\n")); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
		return
	}

	switch method := r.Method; method {

	case http.MethodGet: // process query and return a json
		result, err = getQueryFromDatabase(query, dbConnection, app)
	case http.MethodDelete: // process deletes on the database
		result, err = processDelete(query, dbConnection, app.Log)
	case http.MethodPatch: // process updates on the database
		result, err = processUpdate(query, dbConnection, app.Log)
	case http.MethodPost: // process inserts on the database
		result, err = processInsert(query, dbConnection, app.Log)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("Method not allowed\n")); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
		return
	}

	if err != nil {
		errMessage := fmt.Sprintf("Error processing request %v", err)
		app.Log.Error(errMessage)

		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(errMessage)); err != nil {
			app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		app.Log.Error(fmt.Sprintf("Error writing to buffer %v", err))
	}
}

// processes the update on database
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

// processes the update on database
func processInsert(command string, db *sql.DB, Log *logger.Logger) ([]byte, error) {
	start := time.Now()
	err := validateInsert(command)
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
