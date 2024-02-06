package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

type connectionInfo struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func GetConnectionInfo(cfg *ini.File) string {
	port, _ := cfg.Section("POSTGRES").Key("port").Int()
	connInfo := connectionInfo{
		host:     cfg.Section("POSTGRES").Key("host").String(),
		port:     port,
		user:     cfg.Section("POSTGRES").Key("user").String(),
		password: cfg.Section("POSTGRES").Key("password").String(),
		dbname:   cfg.Section("POSTGRES").Key("dbname").String(),
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connInfo.host,
		connInfo.port,
		connInfo.user,
		connInfo.password,
		connInfo.dbname,
	)
	fmt.Printf(
		"Connecting to postgres host=%v port=%v dbname=%v\n",
		connInfo.host,
		connInfo.port,
		connInfo.dbname,
	)

	return psqlInfo
}

func ConnectToPsql(psqlInfo string) (*sql.DB, error) {
	// return nil, nil
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error: Could not establish a connection with the portgres database", err.Error())
		panic(err)
	}
	log.Println("Successfully connected!")
	return db, nil
}
