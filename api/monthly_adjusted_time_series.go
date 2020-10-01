package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	MONTHLY_ADJUSTED_TIME_SERIES_FUNCTION   = "TIME_SERIES_MONTHLY_ADJUSTED"
	MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME = "MonthlyAdjustedTimeSeries"
)

var (
	MONTHLY_ADJUSTED_TIME_SERIES_HEADERS = []string{
		"id",
		"date",
		"open",
		"high",
		"low",
		"close",
		"adjusted_close",
		"volume",
		"dividend_amount",
	}
	MONTHLY_ADJUSTED_TIME_SERIES_COLUMNS = []string{
		"id INT",
		"date DATE NOT NULL",
		"open FLOAT NOT NULL",
		"high FLOAT NOT NULL",
		"low FLOAT NOT NULL",
		"close FLOAT NOT NULL",
		"adjusted_close FLOAT NOT NULL",
		"volume INT NOT NULL",
		"dividend_amount FLOAT NOT NULL",
		"FOREIGN KEY (id) REFERENCES Ticker(id)",
		"PRIMARY KEY (id, date)",
	}
)

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

type MonthlyAdjustedTimeSeriesService interface {
	Get(symbol string) (MonthlyAdjustedTimeSeries, error)
	Insert(ts MonthlyAdjustedTimeSeries, db database.SQLClient) error
	Sync(symbol string, db database.SQLClient) error
	CreateTable(db database.SQLClient) error
	DropTable(db database.SQLClient) error
}

type monthlyAdjustedTimeSeriesServiceOptions struct {
	Symbol string
}

func newMonthlyAdjustedTimeSeriesServiceOptions(symbol string) monthlyAdjustedTimeSeriesServiceOptions {
	return monthlyAdjustedTimeSeriesServiceOptions{Symbol: symbol}
}

func (o monthlyAdjustedTimeSeriesServiceOptions) ToQueryString() string {
	return fmt.Sprintf("&function=%s&symbol=%s", MONTHLY_ADJUSTED_TIME_SERIES_FUNCTION, o.Symbol)
}

type monthlyAdjustedTimeSeriesServicer struct {
	base baseClient
}

func newMonthlyAdjustedTimeSeriesService(base baseClient) MonthlyAdjustedTimeSeriesService {
	return monthlyAdjustedTimeSeriesServicer{
		base: base,
	}
}

func (s monthlyAdjustedTimeSeriesServicer) Get(symbol string) (MonthlyAdjustedTimeSeries, error) {
	options := newMonthlyAdjustedTimeSeriesServiceOptions(symbol)
	resp, err := s.base.call(options)
	if err != nil {
		return MonthlyAdjustedTimeSeries{}, err
	}

	ts, err := parseMonthlyAdjustedTimeSeries(resp)
	if err != nil {
		return MonthlyAdjustedTimeSeries{}, err
	}

	return ts, nil
}

func (s monthlyAdjustedTimeSeriesServicer) Insert(ts MonthlyAdjustedTimeSeries, db database.SQLClient) error {
	tickerID, err := db.GetTickerID(ts.Metadata.Symbol)
	if err != nil {
		return err
	}
	values := make([]interface{}, 0)
	for _, v := range ts.TimeSeries {
		values = append(values, tickerID)
		values = append(values, v.Date)
		values = append(values, v.Open)
		values = append(values, v.High)
		values = append(values, v.Low)
		values = append(values, v.Close)
		values = append(values, v.AdjustedClose)
		values = append(values, v.Volume)
		values = append(values, v.DividendAmount)
	}
	return nil
}

func (s monthlyAdjustedTimeSeriesServicer) Sync(symbol string, db database.SQLClient) error {
	// TODO
	return nil
}

func (s monthlyAdjustedTimeSeriesServicer) CreateTable(db database.SQLClient) error {
	return db.CreateTable(MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME, MONTHLY_ADJUSTED_TIME_SERIES_COLUMNS)
}

func (s monthlyAdjustedTimeSeriesServicer) DropTable(db database.SQLClient) error {
	return db.DropTable(MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME)
}

func parseMonthlyAdjustedTimeSeries(resp *http.Response) (MonthlyAdjustedTimeSeries, error) {
	target := &MonthlyAdjustedTimeSeriesResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return MonthlyAdjustedTimeSeries{}, err
	}

	timeSeries := target.MonthlyAdjustedTimeSeries

	// slice to hold keys
	keys := make([]string, len(timeSeries))
	i := 0
	for k := range timeSeries {
		keys[i] = k
		i++
	}

	monthlyAdjustedTimeSeriesEntries := make([]MonthlyAdjustedTimeSeriesEntry, len(timeSeries))
	for i, key := range keys {
		entry := timeSeries[key]
		entry.Date = key
		monthlyAdjustedTimeSeriesEntries[i] = entry
	}

	return MonthlyAdjustedTimeSeries{Metadata: target.Metadata, TimeSeries: monthlyAdjustedTimeSeriesEntries}, nil
}
