package main

import (
	"fmt"

	"github.com/mcclurejt/mrkt-backend/api"
)

var API_KEY = "LXCN06KPP1KPOYC2"

func main() {
	s := "AAPL"
	client := api.NewAlphaVantageClient(API_KEY)
	ts, _ := client.MonthlyAdjustedTimeSeriesService.Get(s)

	fmt.Printf("%v", ts)
}
