package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var Metric = "market"
var CryptoTimeSeriesFunction = "price_usd_ohlc"

type CryptoTimeSeries struct {
	TimeSeries []CryptoTimeSeriesEntry
}

type CryptoTimeSeriesResponse struct {
	Timestamp int64                 `json:"t"`
	Data      CryptoTimeSeriesEntry `json:"o"`
}
type CryptoTimeSeriesEntry struct {
	Date  string
	Close float64 `json:"c"`
	High  float64 `json:"h"`
	Low   float64 `json:"l"`
	Open  float64 `json:"o"`
}

func GetCryptoTimeSeries(coin string) CryptoTimeSeries {
	timeSeries := CallCrypto(Metric, CryptoTimeSeriesFunction, coin)
	return timeSeries.(CryptoTimeSeries)
}

func parseCryptoTimeSeries(resp *http.Response) CryptoTimeSeries {
	// [{"t":1600819200,"o":{"c":10493.3637344,"h":10542.9353362,"l":10493.3637344,"o":10541.7558528}}]
	target := &[]CryptoTimeSeriesResponse{}

	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Fatalln(err)
	}

	timeSeries := target

	cryptoTimeSeriesEntries := make([]CryptoTimeSeriesEntry, len(*timeSeries))
	for i, v := range *timeSeries {
		entry := v.Data
		entry.Date = time.Unix(v.Timestamp, 0).Format("2006-01-02")
		cryptoTimeSeriesEntries[i] = entry
	}

	return CryptoTimeSeries{TimeSeries: cryptoTimeSeriesEntries}
}
