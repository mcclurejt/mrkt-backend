package api

import (
	"fmt"
	"log"
	"net/http"
)

var baseURL = "https://www.alphavantage.co/query?"
var apiKey = "LXCN06KPP1KPOYC2"

var cryptoBaseURL = "https://api.glassnode.com/v1/metrics/"
var cryptoApiKey = "105d32cc-afc0-4358-b335-891a35e80736"

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

func CallCrypto(metric string, function string, coin string) interface{} {
	// https://api.glassnode.com/v1/metrics/market/price_usd_close?a=btc&api_key=$API_KEY
	url := fmt.Sprintf(cryptoBaseURL+"%s/"+"%s"+"?a=%s"+"&api_key=%s", metric, function, coin, cryptoApiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	switch function {
	case CryptoTimeSeriesFunction:
		return parseCryptoTimeSeries(resp)
	default:
		return nil
	}

}
