package main

import (
	"fmt"
	"log"

	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/database"
)

var ALPHAVANTAGE_API_KEY = "LXCN06KPP1KPOYC2"
var MARKETSTACK_API_KEY = "02378e09665e4a13b514d5cb29855994"

func main() {
	datasource := "root:1727Clybourn!@tcp(127.0.0.1:3306)/ticker_data"
	db := database.NewMySqlClient(datasource)

	msClient := api.NewMarketStackClient(MARKETSTACK_API_KEY)

	avClient := api.NewAlphaVantageClient(ALPHAVANTAGE_API_KEY)

	err := db.DropAllTables(avClient)
	if err != nil {
		fmt.Println(err.Error())
	}

	// err := db.DropAllTables(msClient)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

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

	err = msClient.TickerService.Insert(tickers, db)
	if err != nil {
		fmt.Println(err.Error())
	}

	rows, err := db.Query("SELECT name FROM Ticker")
	if err != nil {
		log.Fatalln(err.Error())
	}
	var name string
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("Starting Get")
		co, err := avClient.CompanyOverviewService.Get(name)
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println("Starting Insert")
		err = avClient.CompanyOverviewService.Insert(co, db)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}
