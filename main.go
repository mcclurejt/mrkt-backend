package main

import (
	"log"

	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/database"
)

var API_KEY = "LXCN06KPP1KPOYC2"

func main() {
	datasource := "root:1727Clybourn!@tcp(127.0.0.1:3306)/ticker_data"
	db := database.NewMySqlClient(datasource)

	avClient := api.NewAlphaVantageClient(API_KEY)

	err := db.CreateAllTables(avClient)
	if err != nil {
		log.Fatalln(err)
	}

}
