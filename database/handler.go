package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type DatabaseHandler interface {
	// ProcessRequest(w http.ResponseWriter, r *http.Request)
	// ProcessRequestHandlePath(w http.ResponseWriter, r *http.Request)
	// getFromDatabase(query string) ([]byte, error)
	checkDbConnection() error
}

func checkDbConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

// Validate the update command to check for unallowed keywords
func validateUpdate(query string) error {
	words := []string{"delete", "truncate", "drop", "insert", "create "}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

// Validate the insert command to check for unallowed keywords
func validateInsert(query string) error {
	words := []string{"delete", "truncate", "drop", "update"}
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
	words := []string{"update", "truncate", "drop", "insert", "create "}
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
	words := []string{"delete", "truncate", "drop", "update", "insert", "create "}
	lowerQuery := strings.ToLower(query)
	for _, key := range words {
		if strings.Contains(lowerQuery, key) {
			return errors.New(METHOD_NOT_ALLOWED)
		}
	}
	return nil
}

func CloseConn(db *sql.DB) {
	fmt.Println("Closing connection")
	err := db.Close()
	if err != nil {
		fmt.Println("Can't close connection: ", err)
	}
}
