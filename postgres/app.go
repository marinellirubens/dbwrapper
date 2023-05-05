package postgres

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	Db *sql.DB
}

type Teste struct {
	Information string `json:"information" binding:"required"`
	Teste       string `json:"teste" binding:"required"`
}

// example of method to handle a request
func (app *App) GetInfo(c *gin.Context) {
	response := Teste{Information: "OK", Teste: "klsdhkdhfgh"}

	c.Bind(&response)
	fmt.Println(response)

	c.JSON(http.StatusOK, response)
}

func (a *App) GetInfoFromDb(c *gin.Context) {
	query := c.Query("query")
	fmt.Println(query)
	rows, err := a.Db.Query(query)
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
