package main

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/encoder"
)

// nolint
var (
	stop             = make(chan os.Signal, 1)
	exit             = os.Exit
	stderr io.Writer = os.Stderr
)

func main() {
	logger := log.New(stderr, "", log.Lshortfile)

	logger.Println("INITIALIZED")

	tickers := [...]string{"AAPL", "AMZN", "FB", "GOOG", "TSLA"}
	for _, s := range tickers {
		logger.Printf("[%s] : Fetching Data\n", s)
		timeSeries := api.GetMonthlyAdjustedTimeSeries(s)

		logger.Printf("[%s] : Convert to CSV\n", s)
		headers := timeSeries.GetHeaders()
		logger.Printf("[%s] : Header %v\n", s, headers)

		values := timeSeries.GetValues()
		logger.Printf("[%s] : Values:\n%v\n", s, values)

		encoder.CSV(s, headers, values)
		logger.Printf("[%s] : Wrote data to CSV\n", s)
	}

	coins := [...]string{"btc", "eth"}
	for _, c := range coins {
		logger.Printf("[%s] : Fetching Data\n", c)
		cryptoTimeSeries := api.GetCryptoTimeSeries(c)

		logger.Printf("[%s] : Convert to CSV\n", c)
		cryptoHeaders := cryptoTimeSeries.GetHeaders()
		logger.Printf("[%s] : Header : %v\n", c, cryptoHeaders)

		cryptoValues := cryptoTimeSeries.GetValues()
		logger.Printf("[%s] : \nValues : \n %v\n", c, cryptoValues)

		encoder.CSV(strings.ToUpper(c), cryptoHeaders, cryptoValues)
		logger.Printf("[%s] : Wrote data to CSV\n", strings.ToUpper(c))
	}

	logger.Println("COMPLETE")

	exit(0)
}
