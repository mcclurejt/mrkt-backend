package api

import (
	"fmt"
	"log"
	"net/http"
)

var baseURL = "https://www.alphavantage.co/query?"
var apiKey = "LXCN06KPP1KPOYC2"

func Call(function string, symbol string) interface{} {
	url := fmt.Sprintf(baseURL+"function=%s&"+"symbol=%s&"+"apikey=%s", function, symbol, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	switch function {
	case MonthlyAdjustedTimeSeriesFunction:
		return parseMonthlyAdjustedTimeSeries(resp)
	default:
		return nil
	}

}
