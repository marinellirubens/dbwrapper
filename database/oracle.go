package database

import (
	"database/sql"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

func GetOracleConnection(dbInfo DbConnection) *sql.DB {
	connectionString := dbInfo.GetConnString()

	fmt.Println("Opening connection with oracle")
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	return db
}
