package database

import (
	"database/sql"
	"errors"
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

const HAS_TABLE_QUERY = "SHOW TABLES LIKE '?';"

func (m MySqlClient) HasTable(tableName string) bool {
	var existingTable string
	err := m.db.QueryRow(HAS_TABLE_QUERY, tableName).Scan(&existingTable)
	if err != nil {
		return false
	}
	return true
}

const CREATE_TABLE_QUERY = "CREATE TABLE IF NOT EXISTS ?(?);"

func (m MySqlClient) CreateTable(tableName string, columns []string) error {
	_, err := m.db.Exec(CREATE_TABLE_QUERY, tableName, strings.Join(columns, ","))
	if err != nil {
		return err
	}
	return nil
}

const DROP_TABLE_QUERY = "DROP TABLE IF EXISTS ?;"

func (m MySqlClient) DropTable(tableName string) error {
	_, err := m.db.Exec(DROP_TABLE_QUERY, tableName)
	if err != nil {
		return err
	}
	return nil
}

const GET_TICKER_QUERY = "SELECT id FROM Ticker WHERE name = ?;"

func (m MySqlClient) GetTickerID(ticker string) (int, error) {
	if !m.HasTable("Ticker") {
		return -1, errors.New("Ticker table does not exist")
	}

	var tickerID int
	err := m.db.QueryRow(GET_TICKER_QUERY, ticker).Scan(&tickerID)
	if err != nil {
		return -1, err
	}
	return tickerID, nil
}
