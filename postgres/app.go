package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	logs "github.com/marinellirubens/dbwrapper/logger"
)

const METHOD_NOT_ALLOWED = "command not allowed on this endpoint"

// Application object to handle the endpoints and connection with database
type App struct {
	// database connection
	Db *sql.DB
	// logger object for general purposes
	Log *logs.Logger
}

// TODO: treat the path to use the second element as the name of the database
func (app *App) ProcessOracleRequest(w http.ResponseWriter, r *http.Request) {

}

// TODO: treat the path to use the second element as the name of the database
func (app *App) ProcessMongoRequest(w http.ResponseWriter, r *http.Request) {

}

// process selects (GET), delete(DELETE) and update(PATCH)
func (app *App) ProcessPostgresRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	var (
		status int
		result []byte
		err    error
	)
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

// Validate the query command to check for unallowed keywords
func (app *App) validateQuery(query string) error {
	words := []string{"delete", "truncate", "drop", "update"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

// Validate the update command to check for unallowed keywords
func (app *App) validateUpdate(query string) error {
	words := []string{"delete", "truncate", "drop"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

// Validate the delete command to check for unallowed keywords
func (app *App) validateDelete(query string) error {
	words := []string{"update", "truncate", "drop"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

// processes the update on postgresql database
func (app *App) updatePostgres(query string) ([]byte, error) {
	start := time.Now()
	err := app.validateUpdate(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	_, err = app.Db.Exec(query)
	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	return []byte("Success"), nil
}

// process the delete process on postgresql database
func (app *App) deleteFromPostgres(query string) ([]byte, error) {
	start := time.Now()
	err := app.validateDelete(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	_, err = app.Db.Exec(query)
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
	start := time.Now()
	err := app.validateQuery(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}

	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	rows, err := app.Db.Query(query)
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
		for i, _ := range colvals {
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
