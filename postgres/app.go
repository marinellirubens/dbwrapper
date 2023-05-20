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

type App struct {
	// database connection
	Db *sql.DB
	// logger object for general purposes
	Log *logs.Logger
}

func (app *App) ProcessPostgresRequest(w http.ResponseWriter, r *http.Request) {
	var status int
	var result []byte
	var err error

	switch method := r.Method; method {
	case http.MethodGet:
		query := r.URL.Query().Get("query")
		result, err = app.getQueryFromPostgres(query)
	case http.MethodDelete:
		query := r.URL.Query().Get("query")
		result, err = app.updatePostgres(query)
	case http.MethodPatch:
		query := r.URL.Query().Get("query")
		result, err = app.deleteFromPostgres(query)
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

func (app *App) validateQuery(query string) error {
	words := []string{"delete", "truncate", "drop", "update"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New("delete not allowed on this endpoint")
		}
	}
	return nil
}
func (app *App) validateUpdate(query string) error {
	words := []string{"delete", "truncate", "drop"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New("command not allowed on this endpoint")
		}
	}
	return nil
}

func (app *App) validateDelete(query string) error {
	words := []string{"update", "truncate", "drop"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New("command not allowed on this endpoint")
		}
	}
	return nil
}

func (app *App) updatePostgres(query string) ([]byte, error) {
	start := time.Now()
	err := app.validateUpdate(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	return []byte("Success"), nil
}

func (app *App) deleteFromPostgres(query string) ([]byte, error) {
	start := time.Now()
	err := app.validateDelete(query)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err)), err
	}
	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
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
