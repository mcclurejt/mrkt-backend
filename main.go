package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/config"
	"github.com/mcclurejt/mrkt-backend/database"
)

//env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New() //env

	db := database.NewMySqlClient(conf.Db.Datasource)

	msClient := api.NewMarketStackClient(conf.Api.MarketStackAPIKey)

	avClient := api.NewAlphaVantageClient(conf.Api.AlphavantageAPIKey)

	err := db.DropAllTables(avClient)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.DropAllTables(msClient)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.CreateAllTables(msClient)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = db.CreateAllTables(avClient)
	if err != nil {
		fmt.Println(err.Error())
	}

	tickers, err := msClient.TickerService.Get(api.EXCHANGE_NYSE, 20, 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(tickers.Data)

	err = msClient.TickerService.Insert(tickers, db)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, v := range tickers.Data {
		fmt.Println(v)
		ts, err := avClient.MonthlyAdjustedTimeSeriesService.Get(v)
		if err != nil {
			log.Fatalln(err.Error())
		}
		err = avClient.MonthlyAdjustedTimeSeriesService.Insert(ts, db)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}
