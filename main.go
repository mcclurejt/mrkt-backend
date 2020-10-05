package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/config"
	"github.com/mcclurejt/mrkt-backend/database/dynamodb"
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
	avClient := api.NewAlphaVantageClient(conf.Api.AlphavantageAPIKey)
	// gnClient := api.NewGlassNodeClient(conf.Api.GlassNodeAPIKey)
	ddbClient := dynamodb.New()

	err := ddbClient.CreateTable(avClient.CompanyOverviewService)
	if err != nil {
		fmt.Println(err.Error())
	}

	overview, err := avClient.CompanyOverviewService.Get("BABA")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v\n", overview)

	err = ddbClient.PutItem(avClient.CompanyOverviewService, overview)
	if err != nil {
		fmt.Println(err.Error())
	}

}
