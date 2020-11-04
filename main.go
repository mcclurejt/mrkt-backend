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
	util "github.com/mcclurejt/mrkt-backend/api/dynamodbutil"
	iex "github.com/mcclurejt/mrkt-backend/api/iexcloud"
	"github.com/mcclurejt/mrkt-backend/config"
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

	conf := config.New() //env

	// gnClient := gn.NewGlassNodeClient("105d32cc-afc0-4358-b335-891a35e80736")
	// var nupls []*gn.NetUnrealizedProfitLossEntry
	// gnOptions := gn.DefaultNetUnrealizedProfitLossOptions()
	// err := gnClient.BatchCall(gn.NuplRouteName, []string{"BTC", "ETH"}, nupls, gnOptions)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	iexClient := iex.NewIEXCloudClient(conf.Api.IEXCloudAPIKey)
	symbols, _ := iexClient.IexSymbols.Get(context.Background())
	symbol := symbols[0]

	ohlcv, _ := iexClient.Chart.GetSingleDay(context.Background(), symbol.Symbol, "20201103")
	item := ohlcv[0]

	createTableInput, err := util.CreateTableInputFromStruct(iex.OHLCV{})
	if err != nil {
		fmt.Println(err.Error())
	}
	putItemInput, _ := util.PutItemInputFromStruct(item)
	fmt.Println(createTableInput)
	fmt.Println(putItemInput)
	// symbs := []string{"twtr", "amzn"}
	// fmt.Printf("Symbols: %s", strings.Join(symbs, ","))
	// // books, err := iexClient.Book.GetBatch(context.Background(), symbs)
	// // _, err = iexClient.DelayedQuote.Get(context.Background(), "twtr")
	// // _, err = iexClient.IntradayPrices.Get(context.Background(), "twtr")
	// types := []string{"company", "insider-summary", "insider-transactions", "insider-roster"}
	// lt, err := iexClient.Batch.GetSymbolBatch(context.Background(), "amzn", types)
	// fmt.Println(iexSymbs)
	// sp, err := iexClient.SectorPerformance.Get(context.Background())
	// fmt.Println(sp)
	// options := &iex.IntradayOptions{
	// 	ChangeFromClose: true,
	// }
	// _, err = iexClient.IntradayPrices.GetWithOptions(context.Background(), "twtr", options)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// arr := []iex.QueryType{iex.QueryTypeBook, iex.QueryTypeDelayedQuote, iex.QueryTypeCompany}
	// fmt.Println(iex.SliceToString(arr, nil))

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
