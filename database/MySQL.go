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

const HAS_TABLE_QUERY = "SHOW TABLES LIKE '?';"

func (m MySqlClient) HasTable(tableName string) bool {
	var existingTable string
	err := m.db.QueryRow(HAS_TABLE_QUERY, tableName).Scan(&existingTable)
	if err != nil {
		return false
	}
	return true
}

const CREATE_TABLE_QUERY = "CREATE TABLE IF NOT EXISTS ? (?);"

func (m MySqlClient) CreateTable(aps AssetPropertyService) error {
	log.Printf("Creating table '%s' ...", aps.GetTableName())
	stmt, err := m.db.Prepare(CREATE_TABLE_QUERY)
	defer stmt.Close()
	if err != nil {
		return err
	}
	args := []interface{}{aps.GetTableName(), strings.Join(aps.GetTableColumns(), ",")}
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Printf("Failed to create table '%s'.", aps.GetTableName())
		return err
	}
	log.Printf("Table'%s' created", aps.GetTableName())
	return nil
}

func (m MySqlClient) CreateAllTables(client interface{}) error {
	// typ := reflect.TypeOf(client)
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

const DROP_TABLE_QUERY = "DROP TABLE IF EXISTS ?;"

func (m MySqlClient) DropTable(aps AssetPropertyService) error {
	stmt, err := m.db.Prepare(DROP_TABLE_QUERY)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(aps.GetTableName())
	if err != nil {
		return err
	}
	return nil
}

const INSERT_QUERY = "INSERT INTO ? ? VALUES %s"

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

	query := fmt.Sprintf(INSERT_QUERY, strings.Join(valueStrings, ","))
	stmt, err := m.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	// make the array contain {tableName, headerString, values...}
	values = append([]interface{}{headerString}, values)
	values = append([]interface{}{tableName}, values)
	_, err = stmt.Exec(values...)
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
