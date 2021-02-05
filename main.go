package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/joho/godotenv"
	"github.com/mcclurejt/mrkt-backend/api/iex"
	"github.com/mcclurejt/mrkt-backend/config"
	"github.com/sirupsen/logrus"
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

	client := iex.New(
		conf.Api.IEXCloudAPIKey,
		iex.ClientOptionSetTimeout(15*time.Second),
		iex.ClientOptionSetLogLevel(logrus.InfoLevel),
	)
	_, err := client.GetChartSingleDay(context.Background(), "AAPL", time.Now().AddDate(0, 0, -14))
	if err != nil {
		panic(err)
	}
	ohlcvs, err := client.GetChart(context.Background(), "AAPL", iex.ChartRange1m, &iex.ChartOptions{})
	if err != nil {
		panic(err)
	}
	for _, candle := range ohlcvs {
		fmt.Printf("Change: %f, Percent: %f,\n", candle.Change, candle.ChangePercent)
	}

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
