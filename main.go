package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joho/godotenv"
	iex "github.com/mcclurejt/mrkt-backend/api/iexcloud"
)

//env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	// enable cpu profiling
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// conf := config.New() //env

	// msClient := ms.NewMarketStackClient("29250803996c88be8fe2e1ef46dce84e")
	// var tickers []*ms.TickerEntry
	// msOptions := ms.DefaultTickerOptions()
	// msOptions.Limit = 1200
	// err := msClient.BatchCall(ms.TickerRouteName, &tickers, msOptions)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// gnClient := gn.NewGlassNodeClient("105d32cc-afc0-4358-b335-891a35e80736")
	// var nupls []*gn.NetUnrealizedProfitLossEntry
	// gnOptions := gn.DefaultNetUnrealizedProfitLossOptions()
	// err := gnClient.BatchCall(gn.NuplRouteName, []string{"BTC", "ETH"}, nupls, gnOptions)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	iexClient := iex.NewIexCloudClient("pk_1d8a2228abd84b0598a6cf91a5d09f63")
	iexSymbs, err := iexClient.IexSymbols.Get(context.Background())
	// symbs := []string{"twtr", "amzn"}
	// fmt.Printf("Symbols: %s", strings.Join(symbs, ","))
	// // books, err := iexClient.Book.GetBatch(context.Background(), symbs)
	// // _, err = iexClient.DelayedQuote.Get(context.Background(), "twtr")
	// // _, err = iexClient.IntradayPrices.Get(context.Background(), "twtr")
	// types := []string{"company", "insider-summary", "insider-transactions", "insider-roster"}
	// lt, err := iexClient.Batch.GetSymbolBatch(context.Background(), "amzn", types)
	fmt.Println(iexSymbs)
	// sp, err := iexClient.SectorPerformance.Get(context.Background())
	// fmt.Println(sp)
	// options := &iex.IntradayOptions{
	// 	ChangeFromClose: true,
	// }
	// _, err = iexClient.IntradayPrices.GetWithOptions(context.Background(), "twtr", options)
	if err != nil {
		fmt.Println(err.Error())
	}

	// arr := []iex.QueryType{iex.QueryTypeBook, iex.QueryTypeDelayedQuote, iex.QueryTypeCompany}
	// fmt.Println(iex.SliceToString(arr, nil))

	// avClient := av.NewAlphaVantageClient("LXCN06KPP1KPOYC2")
	// var dailyTimeSeries []*av.DailyAdjustedTimeSeriesEntry
	// dailyTimeSeriesOptions := &av.DailyAdjustedTimeSeriesOptions{OutputSize: av.OutputSizeDefault}
	// err = avClient.BatchCall(av.DailyTimeSeriesRouteName, []string{"BABA", "BRK-A", "BRK-B", "AAPL"}, &dailyTimeSeries, dailyTimeSeriesOptions)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// save memory profiling for last
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
