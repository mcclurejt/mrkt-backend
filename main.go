package main

import (
	"fmt"

	"github.com/mcclurejt/mrkt-backend/api"
)

func main() {
	s := "AAPL"
	timeSeries := api.GetMonthlyAdjustedTimeSeries(s)
	fmt.Printf("%v", timeSeries)
}
