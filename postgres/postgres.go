package postgres

import (
	"database/sql"
	"fmt"

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
		connInfo.host, connInfo.port, connInfo.user, connInfo.password, connInfo.dbname)
	//fmt.Println(psqlInfo)

	return psqlInfo
}

func ConnectToPsql(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db, nil
}
