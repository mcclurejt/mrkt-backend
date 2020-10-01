package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlClient struct {
	db *sql.DB
}

func (m MySqlClient) New(datasource string) MySqlClient {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}
	return MySqlClient{db: db}
}

func (m MySqlClient) Insert(tableName string, headers []string, values []interface{}) error {
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
	stmt := fmt.Sprintf("INSERT INTO %s %s VALUES %s ", tableName, headerString, strings.Join(valueStrings, ","))
	_, err := m.db.Exec(stmt, values...)
	return err
}

func (m MySqlClient) CreateTable(tableName string, columns []string) error {
	query := "CREATE TABLE IF NOT EXISTS ?(?);"
	_, err := m.db.Exec(query, tableName, strings.Join(columns, ","))
	if err != nil {
		return err
	}
	return nil
}

func (m MySqlClient) DropTable(tableName string) error {
	query := "DROP TABLE IF EXISTS ?;"
	_, err := m.db.Exec(query, tableName)
	if err != nil {
		return err
	}
	return nil
}
