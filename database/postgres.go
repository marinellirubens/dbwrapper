package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/marinellirubens/dbwrapper/internal/logger"
)

func GetPostgresConnection(psqlInfo DbConnection, log *logger.Logger) (*sql.DB, error) {
	log.Debug("Opening connection with oracle")

	db, err := sql.Open("postgres", psqlInfo.GetConnString())
	if err != nil {
		log.Error(fmt.Sprintf("Error connecting to postgres %v\n", err))
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("Error: Could not establish a connection with the portgres database %s", err.Error()))
		return db, err
	}

	log.Info("Successfully connected!")
	return db, nil
}
