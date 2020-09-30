package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Store interface {
	SetupDatabase(d string)
	CreateDatabase(d string)
	ChooseDatabase(d string)
	CreateTable(tbl string, columns string)
	CloseDatabase()
	ResetTables()
	DropTable(tbl string)
	LoadCsv(ticker string, tbl string)
}

func (store *dbStore) SetupDatabase(d string) {
	store.CreateDatabase(d)
	store.ChooseDatabase(d)
	store.ResetTables()

	var (
		ticker  = "Ticker"
		columns = "id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR (8) NOT NULL UNIQUE"
	)
	store.CreateTable(ticker, columns)

	ticker = "Candle"
	columns = "id INT, date DATE NOT NULL, open FLOAT NOT NULL, high FLOAT NOT NULL, low FLOAT NOT NULL, close FLOAT NOT NULL, FOREIGN KEY (id) REFERENCES Ticker(id), PRIMARY KEY (id, date)"
	store.CreateTable(ticker, columns)
}

func (store *dbStore) CreateDatabase(d string) {
	_, err := store.db.Exec("CREATE DATABASE IF NOT EXISTS ticker_data")
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Successfully created database ticker_data")
	}
}

func (store *dbStore) ChooseDatabase(d string) {
	_, err := store.db.Exec("USE " + d)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("ticker_data (DB) selected successfully")
	}
}

func (store *dbStore) ResetTables() {
	store.DropTable("Candle")
	store.DropTable("Ticker")
}

func (store *dbStore) DropTable(tbl string) {
	_, err := store.db.Exec("DROP TABLE IF EXISTS " + tbl)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Successfully dropped table %s\n", tbl)
	}
}

func (store *dbStore) CreateTable(tbl string, columns string) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s);", tbl, columns)
	stmt, err := store.db.Prepare(query)
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Table, %s, created successfully\n", tbl)
	}
}

func (store *dbStore) LoadCsv(ticker string, tbl string) {
	mysql.RegisterLocalFile("data/time_series/" + ticker + ".csv")
	query := fmt.Sprintf("LOAD DATA LOCAL INFILE 'data/time_series/%s.csv' INTO TABLE %s FIELDS TERMINATED BY ',' LINES TERMINATED BY ';\n' IGNORE 1 LINES (Date, Open, High, Low, Close) SET id = (SELECT id FROM Ticker WHERE name = \"%s\"), date = STR_TO_DATE(@Date, '%%Y-%%m-%%d');", ticker, tbl, ticker)
	_, err := store.db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Loaded %s time series into %s successfully\n", ticker, tbl)
	}
}
func (store *dbStore) CloseDatabase() {
	defer store.db.Close()
}

var store Store

func InitStore(s Store) {
	store = s
}
