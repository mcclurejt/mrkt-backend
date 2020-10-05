package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mcclurejt/mrkt-backend/api/alphavantage"
	"github.com/mcclurejt/mrkt-backend/config"
)

//env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New() //env

	// db := database.NewMySqlClient(conf.Db.Datasource)

	// msClient := api.NewMarketStackClient(conf.Api.MarketStackAPIKey)
	avClient := alphavantage.NewAlphaVantageClient(conf.Api.AlphavantageAPIKey)
	// gnClient := api.NewGlassNodeClient(conf.Api.GlassNodeAPIKey)
	// ddbClient := dynamodb.New()

	// err := ddbClient.CreateTable(avClient.DailyAdjustedTimeSeriesService)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	series, err := avClient.DailyAdjustedTimeSeries.Get("BABA", alphavantage.OutputSizeDefault)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, entry := range series.TimeSeries {
		fmt.Printf("%v\n", entry)
	}

	// err = ddbClient.PutAllItems(avClient.DailyAdjustedTimeSeriesService, series.TimeSeries)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

}
