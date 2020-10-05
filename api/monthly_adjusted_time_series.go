package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mcclurejt/mrkt-backend/database"
	db "github.com/mcclurejt/mrkt-backend/database/dynamodb"
)

const (
	MONTHLY_ADJUSTED_TIME_SERIES_FUNCTION   = "TIME_SERIES_MONTHLY_ADJUSTED"
	MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME = "MonthlyAdjustedTimeSeries"
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
	Symbol         string
	Open           float64 `json:"1. open,string"`
	High           float64 `json:"2. high,string"`
	Low            float64 `json:"3. low,string"`
	Close          float64 `json:"4. close,string"`
	AdjustedClose  float64 `json:"5. adjusted close,string"`
	Volume         int     `json:"6. volume,string"`
	DividendAmount float64 `json:"7. open,string"`
}

type MonthlyAdjustedTimeSeriesService interface {
	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput

	Get(symbol string) (MonthlyAdjustedTimeSeries, error)
	Sync(symbol string, db database.SQLClient) error
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

func (s monthlyAdjustedTimeSeriesServicer) GetCreateTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Symbol"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Date"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Symbol"),
				KeyType:       aws.String("RANGE"),
			},
		},
		BillingMode: aws.String(db.DefaultBillingMode),
		TableName:   aws.String(MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME),
	}
}

func (s monthlyAdjustedTimeSeriesServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(MONTHLY_ADJUSTED_TIME_SERIES_TABLE_NAME),
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
		_, ok := err.(*AlphaVantageRateExceededError)
		if ok {
			time.Sleep(DEFAULT_RETRY_PERIOD_SECONDS * time.Second)
			return s.Get(symbol)
		}
		return MonthlyAdjustedTimeSeries{}, err
	}

	return ts, nil
}

func (s monthlyAdjustedTimeSeriesServicer) Sync(symbol string, db database.SQLClient) error {
	// TODO
	return nil
}

func parseMonthlyAdjustedTimeSeries(resp *http.Response) (MonthlyAdjustedTimeSeries, error) {
	target := &MonthlyAdjustedTimeSeriesResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return MonthlyAdjustedTimeSeries{}, err
	}

	timeSeries := target.MonthlyAdjustedTimeSeries

	// check to see if the rate was exceeded and no objects were returned (still gives 200 status code)
	if len(timeSeries) < 1 {
		return MonthlyAdjustedTimeSeries{}, &AlphaVantageRateExceededError{}
	}

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
		entry.Symbol = target.Metadata.Symbol
		monthlyAdjustedTimeSeriesEntries[i] = entry
	}

	return MonthlyAdjustedTimeSeries{Metadata: target.Metadata, TimeSeries: monthlyAdjustedTimeSeriesEntries}, nil
}
