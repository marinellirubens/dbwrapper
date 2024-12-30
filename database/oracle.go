package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/sijms/go-ora/v2"
)

func GetSqlDBWithPureDriver(dbParams map[string]string) *sql.DB {
	connectionString := fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		dbParams["username"],
		dbParams["password"],
		dbParams["server"],
		dbParams["port"],
		dbParams["service"],
	)
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

var localDB = map[string]string{
	"service":  "LGBRTMST",
	"username": "tms_if",
	"server":   "136.166.34.123",
	"port":     "3006",
	"password": "Qtms108!",
}

func sqlOperations(db *sql.DB) {
	var queryResultColumnOne, queryResultColumnTwo string
	fmt.Println("collecting info from oracle")
	rows, err := db.Query("SELECT to_char(systimestamp,'HH24:MI:SS') as time, 'text' as message FROM dual")
	cols, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	allgeneric := make([]map[string]interface{}, 0)
	colvals := make([]interface{}, len(cols))
	for rows.Next() {
		colassoc := make(map[string]interface{}, len(cols))
		for i := range colvals {
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

	js, _ := json.Marshal(allgeneric) //, "", "    ")
	fmt.Printf("%s\n", js)
	fmt.Println("The time in the database ", queryResultColumnOne, queryResultColumnTwo)
}

func closeConn(db *sql.DB) {
	fmt.Println("Closing connection")
	err := db.Close()
	if err != nil {
		fmt.Println("Can't close connection: ", err)
	}
}

//func main() {
//db := GetSqlDBWithPureDriver(localDB)
//defer closeConn(db)
//sqlOperations(db)
//}
