package database

import (
	"database/sql"
	"fmt"

	"github.com/marinellirubens/dbwrapper/internal/logger"
	_ "github.com/sijms/go-ora/v2"
)

func GetOracleConnection(dbInfo DbConnection, log *logger.Logger) (*sql.DB, error) {
	log.Debug("Opening connection with oracle")

	db, err := sql.Open("oracle", dbInfo.GetConnString())
	if err != nil {
		log.Error(fmt.Sprintf("Error connecting to oracle %v\n", err))
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("Error: Could not establish a connection with the oracle database %s", err.Error()))
		return db, err
	}

	log.Info("Successfully connected!")
	return db, nil
}
