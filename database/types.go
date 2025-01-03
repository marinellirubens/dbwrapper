package database

import (
	"encoding/json"
	"fmt"

	"github.com/marinellirubens/dbwrapper/internal/config"
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
}

// implements [DbConnection] interface
type OracleConnectionInfo struct {
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

var dbConnections map[string]DbConnection = map[string]DbConnection{}

func SetOracleConnection(connInfo config.Database) OracleConnectionInfo {
	// Add an OracleConnectionInfo
	oracleConn := OracleConnectionInfo{
		Server:   connInfo.Host,
		Port:     connInfo.Port,
		User:     connInfo.User,
		password: connInfo.Password,
		Service:  connInfo.Service,
		Dbtype:   ORACLE,
	}
	return oracleConn
}

func SetPostgresConnection(connInfo config.Database) PgConnectionInfo {
	// Add an OracleConnectionInfo
	pgConn := PgConnectionInfo{
		Host:     connInfo.Host,
		Port:     connInfo.Port,
		User:     connInfo.User,
		password: connInfo.Password,
		Dbname:   connInfo.Dbname,
		Dbtype:   POSTGRES,
	}

	return pgConn
}

func SetMapping() string {
	// Add an OracleConnectionInfo
	oracleConn := OracleConnectionInfo{
		Server:   "oracle.example.com",
		Port:     1521,
		User:     "oracle_user",
		password: "secure_password",
		Service:  "ORCL",
		Dbtype:   ORACLE,
	}
	dbConnections["orcl"] = &oracleConn

	//fmt.Printf("%s\n", js)

	// Add a PgConnectionInfo
	pgConn := &PgConnectionInfo{
		Host:     "postgres.example.com",
		Port:     5432,
		User:     "pg_user",
		password: "secure_password",
		Dbname:   "example_db",
		Dbtype:   POSTGRES,
	}
	dbConnections["localdb"] = pgConn

	//info := dbConnections["oracle"]
	//fmt.Printf("printing connection %v", info.GetConnInfo())
	js, err := json.MarshalIndent(dbConnections, "", "    ")
	if err != nil {
		fmt.Println("error")
		panic(err)
	}
	fmt.Printf("%s\n", js)

	return ""
}
