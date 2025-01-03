package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

func GetConnectionInfo(cfg *ini.File) string {
	port, _ := cfg.Section("POSTGRES").Key("port").Int()
	connInfo := PgConnectionInfo{
		Host:     cfg.Section("POSTGRES").Key("host").String(),
		Port:     port,
		User:     cfg.Section("POSTGRES").Key("user").String(),
		password: cfg.Section("POSTGRES").Key("password").String(),
		Dbname:   cfg.Section("POSTGRES").Key("dbname").String(),
	}

	connString := connInfo.GetConnString()
	fmt.Println(connInfo.GetConnInfo())

	return connString
}

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
