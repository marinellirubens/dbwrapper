package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

type app struct {
	db *sql.DB
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
	fmt.Println(psqlInfo)

	return psqlInfo
}

func ConnectToPsql(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db, nil
}

func (a *app) GetInfoFromDb(c *gin.Context) {
	query := c.Query("query")
	fmt.Println(query)
	rows, err := a.db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, "error")
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
	// fmt.Println(allgeneric)
	c.JSON(200, allgeneric)
}
