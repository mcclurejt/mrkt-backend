package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

var MonthlyAdjustedTimeSeriesFunction = "TIME_SERIES_MONTHLY_ADJUSTED"

type MonthlyAdjustedTimeSeries struct {
	Metadata   MonthlyAdjustedTimeSeriesMetadata
	TimeSeries []MonthlyAdjustedTimeSeriesEntry
}

type MonthlyAdjustedTimeSeriesResponse struct {
	Metadata                  MonthlyAdjustedTimeSeriesMetadata         `json:"Meta Data"`
	MonthlyAdjustedTimeSeries map[string]MonthlyAdjustedTimeSeriesEntry `json:"Monthly Adjusted Time Series"`
}

type MonthlyAdjustedTimeSeriesMetadata struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	TimeZone      string `json:"4. Time Zone"`
}

type MonthlyAdjustedTimeSeriesEntry struct {
	Date           string
	Open           float64 `json:"1. open,string"`
	High           float64 `json:"2. high,string"`
	Low            float64 `json:"3. low,string"`
	Close          float64 `json:"4. close,string"`
	AdjustedClose  float64 `json:"5. adjusted close,string"`
	Volume         int     `json:"6. volume,string"`
	DividendAmount float64 `json:"7. open,string"`
}

func GetMonthlyAdjustedTimeSeries(symbol string) MonthlyAdjustedTimeSeries {
	timeSeries := Call(MonthlyAdjustedTimeSeriesFunction, symbol)
	return timeSeries.(MonthlyAdjustedTimeSeries)
}

func parseMonthlyAdjustedTimeSeries(resp *http.Response) MonthlyAdjustedTimeSeries {
	target := &MonthlyAdjustedTimeSeriesResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		log.Fatalln(err)
	}

	timeSeries := target.MonthlyAdjustedTimeSeries

	// slice to hold keys
	keys := make([]string, len(timeSeries))

	i := 0
	for k, _ := range timeSeries {
		keys[i] = k
		i++
	}

	monthlyAdjustedTimeSeriesEntries := make([]MonthlyAdjustedTimeSeriesEntry, len(timeSeries))
	for i, key := range keys {
		entry := timeSeries[key]
		entry.Date = key
		monthlyAdjustedTimeSeriesEntries[i] = entry
	}

	return MonthlyAdjustedTimeSeries{Metadata: target.Metadata, TimeSeries: monthlyAdjustedTimeSeriesEntries}

}

func (m MonthlyAdjustedTimeSeries) GetHeaders() []string {
	e := MonthlyAdjustedTimeSeriesEntry{}  // type
	t := reflect.ValueOf(&e).Elem().Type() // easily abstracted out -> START
	r := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		r[i] = t.Field(i).Name
	}
	return r // easily abstracted out -> END
}

func (m MonthlyAdjustedTimeSeries) GetValues() [][]string {
	var r [][]string // easily abstracted out -> START
	fmt.Printf("%v", m)
	for _, v := range m.TimeSeries {
		item := reflect.ValueOf(v)
		var record []string
		for i := 0; i < item.NumField(); i++ {
			itm := item.Field(i).Interface()
			record = append(record, fmt.Sprintf("%v", itm))
		}
		r = append(r, record)
	}
	return r // easily abstracted out -> END
}
