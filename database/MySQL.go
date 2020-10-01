package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlClient struct {
	db *sql.DB
}

func NewMySqlClient(datasource string) MySqlClient {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
	}
	return MySqlClient{db: db}
}

const HAS_TABLE_QUERY = "SHOW TABLES LIKE '%s';"

func (m MySqlClient) HasTable(tableName string) bool {
	query := fmt.Sprintf(HAS_TABLE_QUERY, tableName)
	var existingTable string
	err := m.db.QueryRow(query).Scan(&existingTable)
	if err != nil {
		return false
	}
	return true
}

const CREATE_TABLE_QUERY = "CREATE TABLE IF NOT EXISTS %s ( %s )"

func (m MySqlClient) CreateTable(aps AssetPropertyService) error {
	log.Printf("Creating table '%s' ...", aps.GetTableName())

	query := fmt.Sprintf(CREATE_TABLE_QUERY, aps.GetTableName(), strings.Join(aps.GetTableColumns(), ","))
	_, err := m.db.Exec(query)
	if err != nil {
		log.Printf("Failed to create table '%s'.", aps.GetTableName())
		return err
	}
	log.Printf("Table'%s' created", aps.GetTableName())
	return nil
}

func (m MySqlClient) CreateAllTables(client interface{}) error {
	val := reflect.ValueOf(client)
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.CanInterface() {
			service, ok := f.Interface().(AssetPropertyService)
			if ok {
				err := m.CreateTable(service)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

const DROP_TABLE_QUERY = "DROP TABLE IF EXISTS %s;"

func (m MySqlClient) DropTable(aps AssetPropertyService) error {
	log.Printf("Dropping table '%s' ...", aps.GetTableName())

	query := fmt.Sprintf(DROP_TABLE_QUERY, aps.GetTableName())
	_, err := m.db.Exec(query)
	if err != nil {
		log.Printf("Failed to drop table '%s'.", aps.GetTableName())
		return err
	}
	log.Printf("Table'%s' dropped", aps.GetTableName())
	return nil
}

func (m MySqlClient) DropAllTables(client interface{}) error {
	val := reflect.ValueOf(client)
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.CanInterface() {
			service, ok := f.Interface().(AssetPropertyService)
			if ok {
				err := m.DropTable(service)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

const INSERT_QUERY = "INSERT INTO %s %s VALUES %s"

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

	query := fmt.Sprintf(INSERT_QUERY, tableName, headerString, strings.Join(valueStrings, ","))
	_, err := m.db.Exec(query, values...)
	if err != nil {
		return err
	}
	log.Printf("Inserted values into table '%s'", tableName)
	return nil
}

const GET_TICKER_QUERY = `SELECT id FROM Ticker WHERE name="%s"`

func (m MySqlClient) GetTickerID(ticker string) (int, error) {
	if !m.HasTable("Ticker") {
		return -1, errors.New("Ticker table does not exist")
	}

	query := fmt.Sprintf(GET_TICKER_QUERY, ticker)
	var tickerID int
	err := m.db.QueryRow(query).Scan(&tickerID)
	if err != nil {
		return -1, err
	}
	return tickerID, nil
}

func (m MySqlClient) Query(query string) (*sql.Rows, error) {
	return m.db.Query(query)
}
