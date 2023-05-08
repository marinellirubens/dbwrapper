package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	logs "github.com/marinellirubens/dbwrapper/logger"
	//"github.com/gin-gonic/gin"
)

type App struct {
	// database connection
	Db *sql.DB
	// logger object for general purposes
	Log *logs.Logger
}

type Teste struct {
	Information string `json:"information" binding:"required"`
	Teste       string `json:"teste" binding:"required"`
}

func (app *App) GetInfoNative(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v %v %v %v\n", r.Method, r.URL.Path, r.Host, r.Proto)
	js, _ := json.Marshal(Teste{Information: "OK", Teste: "klsdhkdhfgh"})

	w.WriteHeader(203)
	w.Write(js)
}

func (app *App) GetInfoNativeTeste(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GetInfoTeste %v %v %v %v\n", r.Method, r.URL.Path, r.Host, r.Proto)
	w.Write([]byte("teste"))
}

func (app *App) GetInfoFromDb(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := r.URL.Query().Get("query")
	app.Log.Info(fmt.Sprintf("Query sent: `%s` processing...", query))
	rows, err := app.Db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("%v", err)))
	}
	//fmt.Println(rows)
	//fmt.Println(err)
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

	app.Log.Debug(fmt.Sprintf("Processed in %vus", time.Since(start).Microseconds()))
	js, _ := json.Marshal(allgeneric)
	w.WriteHeader(203)
	w.Write(js)

}
