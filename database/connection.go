package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/marinellirubens/dbwrapper/internal/config"
	"github.com/marinellirubens/dbwrapper/internal/logger"
	_ "github.com/sijms/go-ora/v2"
)

func SetOracleConnection(connInfo config.Database) OracleConnectionInfo {
	// Add an OracleConnectionInfo
	oracleConn := OracleConnectionInfo{
		Id:       connInfo.Id,
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
		Id:       connInfo.Id,
		Host:     connInfo.Host,
		Port:     connInfo.Port,
		User:     connInfo.User,
		password: connInfo.Password,
		Dbname:   connInfo.Dbname,
		Dbtype:   POSTGRES,
	}

	return pgConn
}

func SetMysqlConnection(connInfo config.Database) MysqlConnectionInfo {
	// Add an OracleConnectionInfo
	pgConn := MysqlConnectionInfo{
		Id:       connInfo.Id,
		Host:     connInfo.Host,
		Port:     connInfo.Port,
		User:     connInfo.User,
		password: connInfo.Password,
		Dbname:   connInfo.Dbname,
		Dbtype:   MYSQL,
	}

	return pgConn
}

// Returns a connection for the database itentifiying using the dbInfo to get the connection string and database type
func GetConnection(dbInfo DbConnection, log *logger.Logger) (*sql.DB, error) {
	log.Debug(fmt.Sprintf("Opening connection with %s", dbInfo.GetDbType()))

	db, err := sql.Open(dbInfo.GetDbType(), dbInfo.GetConnString())
	if err != nil {
		log.Error(fmt.Sprintf("Error connecting to %s %v\n", dbInfo.GetDbType(), err))
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("Error: Could not establish a connection with %s database %s", dbInfo.GetDbType(), err.Error()))
		return db, err
	}

	log.Info("Successfully connected!")
	return db, nil
}
