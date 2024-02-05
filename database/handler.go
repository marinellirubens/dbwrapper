package database

import (
	"database/sql"
	"errors"
	"strings"
)

type DatabaseHandler interface {
	// ProcessRequest(w http.ResponseWriter, r *http.Request)
	// ProcessRequestHandlePath(w http.ResponseWriter, r *http.Request)
	// getFromDatabase(query string) ([]byte, error)
	checkDbConnection() error
}

type PostgresHandler struct {
	db                *sql.DB
	connection_string string
}

type OracleHandler struct {
	db                *sql.DB
	connection_string string
}

type MongoHandler struct {
	db                *sql.DB
	connection_string string
}

func (handler PostgresHandler) checkDbConnection() error {
	if err := handler.db.Ping(); err != nil {
		db, err := ConnectToPsql(handler.connection_string)
		if err != nil {
			return err
		}
		handler.db = db
	}
	return nil
}

// Validate the update command to check for unallowed keywords
func validateUpdate(query string) error {
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
func validateDelete(query string) error {
	words := []string{"update", "truncate", "drop"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

// Validate the query command to check for unallowed keywords
func validateQuery(query string) error {
	words := []string{"delete", "truncate", "drop", "update"}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}
