package main

import (
	"fmt"

	"github.com/mcclurejt/mrkt-backend/api"
	"github.com/mcclurejt/mrkt-backend/encoder"
)

func main() {
	tickers := [...]string{"AAPL", "AMZN", "FB", "GOOG", "TSLA"}
	for _, s := range tickers {
		fmt.Printf("[%s] : Fetching Data\n", s)
		timeSeries := api.GetMonthlyAdjustedTimeSeries(s)

		fmt.Printf("[%s] : Convert to CSV\n", s)
		headers := timeSeries.GetHeaders()
		fmt.Printf("[%s] : Header %v\n", s, headers)

		values := timeSeries.GetValues()
		fmt.Printf("[%s] : \nValues \n", s, values)

		encoder.CSV(s, headers, values)
	}

	c := "btc"
	cryptoTimeSeries := api.GetCryptoTimeSeries(c)
	fmt.Printf("%v", cryptoTimeSeries)
	cryptoHeader := cryptoTimeSeries.GetHeaders()
	cryptoValues := cryptoTimeSeries.GetValues()

	encoder.CSV("BTC", cryptoHeader, cryptoValues)

}
