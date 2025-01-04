package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func GetPostgresConnection(psqlInfo DbConnection) (*sql.DB, error) {
	// return nil, nil
	db, err := sql.Open("postgres", psqlInfo.GetConnString())
	if err != nil {
		return db, err
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error: Could not establish a connection with the portgres database", err.Error())
		return db, err
	}
	log.Println("Successfully connected!")
	return db, nil
}
