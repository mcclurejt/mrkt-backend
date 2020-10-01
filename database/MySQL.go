package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySql struct {
	Db *sql.DB
}

func New(datasource string) MySql {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}
	return MySql{Db: db}
}

func (m MySql) Insert(tableName string, headers []string, values []interface{}) error {
	fmt.Printf("Insert")
	rowWidth := len(headers)
	valueStrings := make([]string, 0, len(values)/rowWidth)
	for i := 0; i < len(values); i += rowWidth {
		valueString := "("
		for j := 0; j < rowWidth; j++ {
			if j == rowWidth-1 {
				valueString += "?"
			} else {
				valueString += "?, "
			}
		}
		valueString += ")"
		valueStrings = append(valueStrings, valueString)
	}
	headerString := "(" + strings.Join(headers, ",") + ")"
	stmt := fmt.Sprintf("INSERT INTO %s %s VALUES(%s)", tableName, headerString, strings.Join(valueStrings, ","))
	_, err := m.Db.Exec(stmt, values...)
	return err
}
