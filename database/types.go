package database

import (
	"fmt"
)

const (
	ORACLE   = "oracle"
	POSTGRES = "postgres"
)

type DbConnection interface {
	// returns the connection string using the necessary format for the connection with the database
	GetConnString() string
	// Returns basic connection info to be printed or logged without sensitive info
	GetConnInfo() string
	// Returns the db connection
	GetDbType() string
	// Returns the db name
	GetDbId() string
}

// implements [DbConnection] interface
type OracleConnectionInfo struct {
	Id      string `json:"id"`
	Server  string `json:"server"`
	Port    int    `json:"port"`
	User    string `json:"user"`
	Service string `json:"service"`
	Dbtype  string `json:"dbtype"`

	password string
}

func (conn *OracleConnectionInfo) GetDbType() string {
	return conn.Dbtype
}

func (conn *OracleConnectionInfo) GetDbId() string {
	return conn.Id
}

// returns the connection string using the necessary format for the connection with the database
func (conn *OracleConnectionInfo) GetConnString() string {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		conn.Dbtype,
		conn.User,
		conn.password,
		conn.Server,
		conn.Port,
		conn.Service,
	)

	return connectionString
}

// Returns basic connection info to be printed or logged without sensitive info
func (conn *OracleConnectionInfo) GetConnInfo() string {
	connectionString := fmt.Sprintf("%s://%s@%s:%d/%s",
		conn.Dbtype,
		conn.User,
		conn.Server,
		conn.Port,
		conn.Service,
	)

	return connectionString
}

// implements [DbConnection] interface
type PgConnectionInfo struct {
	Id     string `json:"id"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   string `json:"user"`
	Dbname string `json:"dbname"`
	Dbtype string `json:"dbtype"`

	password string
}

func (conn *PgConnectionInfo) GetDbType() string {
	return conn.Dbtype
}

func (conn *PgConnectionInfo) GetDbId() string {
	return conn.Id
}

// returns the connection string using the necessary format for the connection with the database
func (conn *PgConnectionInfo) GetConnString() string {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conn.Host,
		conn.Port,
		conn.User,
		conn.password,
		conn.Dbname,
	)

	return psqlInfo
}

// Returns basic connection info to be printed or logged without sensitive info
func (conn *PgConnectionInfo) GetConnInfo() string {
	plsqlInfo := fmt.Sprintf(
		"Connecting to %s host=%v port=%v dbname=%v\n",
		conn.Dbtype,
		conn.Host,
		conn.Port,
		conn.Dbname,
	)

	return plsqlInfo
}
