package main

import (
	"fmt"
	"log"

	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/database"
)

var API_KEY = "LXCN06KPP1KPOYC2"

func main() {
	s := "AAPL"
	datasource := "root:1727Clybourn!@tcp(127.0.0.1:3306)/ticker_data"
	db := database.MySql.New(datasource)

	client := api.NewAlphaVantageClient(API_KEY)
	ts, _ := client.MonthlyAdjustedTimeSeriesService.Get(s)
	headers := []string{"name", "date", "open", "high", "low", "close"}
	values := make([]interface{}, 0)
	for _, v := range ts.TimeSeries {
		values = append(values, ts.Metadata.Symbol)
		values = append(values, v.Date)
		values = append(values, v.Open)
		values = append(values, v.High)
		values = append(values, v.Low)
		values = append(values, v.Close)
	}
	fmt.Println(values)

	res, _ := db.Db.Query("SHOW Tables;")
	for res.Next() {
		var tableName string
		res.Scan(&tableName)

		fmt.Println(tableName)
	}
	_, err := db.Db.Exec("USE ticker_data;")
	_, err = db.Db.Exec("CREATE TABLE IF NOT EXISTS MonthlyAdjustedTimeSeries(name VARCHAR(8) NOT NULL, date DATE NOT NULL, open FLOAT NOT NULL, high FLOAT NOT NULL, low FLOAT NOT NULL, close FLOAT NOT NULL, PRIMARY KEY (name, date));")

	err = db.Insert("MonthlyAdjustedTimeSeries", headers, values)
	if err != nil {
		log.Fatal(err)
	}
}
