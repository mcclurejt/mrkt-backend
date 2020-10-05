package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/mcclurejt/mrkt-backend/api/dynamodb"
	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	DAILY_ADJUSTED_TIME_SERIES_FUNCTION   = "TIME_SERIES_DAILY_ADJUSTED"
	DAILY_ADJUSTED_TIME_SERIES_TABLE_NAME = "DailyAdjustedTimeSeries"
	OutputSizeCompact                     = "compact"
	OutputSizeFull                        = "full"
	OutputSizeDefault                     = OutputSizeCompact
)

type DailyAdjustedTimeSeries struct {
	Metadata   *DailyAdjustedTimeSeriesMetadata
	TimeSeries []*DailyAdjustedTimeSeriesEntry
}

type DailyAdjustedTimeSeriesResponse struct {
	Metadata                *DailyAdjustedTimeSeriesMetadata         `json:"Meta Data"`
	DailyAdjustedTimeSeries map[string]*DailyAdjustedTimeSeriesEntry `json:"Time Series (Daily)"`
}

type DailyAdjustedTimeSeriesMetadata struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Last Refreshed"`
}

type DailyAdjustedTimeSeriesEntry struct {
	Date             string
	Symbol           string
	Open             float64 `json:"1. open,string"`
	High             float64 `json:"2. high,string"`
	Low              float64 `json:"3. low,string"`
	Close            float64 `json:"4. close,string"`
	AdjustedClose    float64 `json:"5. adjusted close,string"`
	Volume           int     `json:"6. volume,string"`
	DividendAmount   float64 `json:"7. open,string"`
	SplitCoefficient float64 `json:"8. split coefficient,string"`
}

type DailyAdjustedTimeSeriesService interface {
	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput

	Get(symbol string, outputSize string) (*DailyAdjustedTimeSeries, error)
	Sync(symbol string, db database.SQLClient) error
}

type dailyAdjustedTimeSeriesServiceOptions struct {
	Symbol     string
	OutputSize string
}

func newDailyAdjustedTimeSeriesServiceOptions(symbol string, outputSize string) dailyAdjustedTimeSeriesServiceOptions {
	return dailyAdjustedTimeSeriesServiceOptions{Symbol: symbol, OutputSize: outputSize}
}

func (o dailyAdjustedTimeSeriesServiceOptions) ToQueryString() string {
	return fmt.Sprintf("&function=%s&symbol=%s&outputsize=%s", DAILY_ADJUSTED_TIME_SERIES_FUNCTION, o.Symbol, o.OutputSize)
}

type dailyAdjustedTimeSeriesServicer struct {
	base *baseClient
}

func newDailyAdjustedTimeSeriesService(base *baseClient) DailyAdjustedTimeSeriesService {
	return dailyAdjustedTimeSeriesServicer{
		base: base,
	}
}

func (s dailyAdjustedTimeSeriesServicer) GetCreateTableInput() *dynamodb.CreateTableInput {
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
		TableName:   aws.String(DAILY_ADJUSTED_TIME_SERIES_TABLE_NAME),
	}
}

func (s dailyAdjustedTimeSeriesServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(DAILY_ADJUSTED_TIME_SERIES_TABLE_NAME),
	}
}

func (s dailyAdjustedTimeSeriesServicer) Get(symbol string, outputSize string) (*DailyAdjustedTimeSeries, error) {
	options := newDailyAdjustedTimeSeriesServiceOptions(symbol, outputSize)
	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	ts, err := parseDailyAdjustedTimeSeries(resp)
	if err != nil {
		_, ok := err.(*AlphaVantageRateExceededError)
		if ok {
			time.Sleep(defaultRetryPeriod * time.Second)
			return s.Get(symbol, outputSize)
		}
		return nil, err
	}

	return ts, nil
}

func (s dailyAdjustedTimeSeriesServicer) Sync(symbol string, db database.SQLClient) error {
	// TODO
	return nil
}

func parseDailyAdjustedTimeSeries(resp *http.Response) (*DailyAdjustedTimeSeries, error) {
	target := &DailyAdjustedTimeSeriesResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	timeSeries := target.DailyAdjustedTimeSeries

	// check to see if the rate was exceeded and no objects were returned (still gives 200 status code)
	if len(timeSeries) < 1 {
		return nil, &AlphaVantageRateExceededError{}
	}

	// slice to hold keys
	keys := make([]string, len(timeSeries))
	i := 0
	for k := range timeSeries {
		keys[i] = k
		i++
	}

	dailyAdjustedTimeSeriesEntries := make([]*DailyAdjustedTimeSeriesEntry, len(timeSeries))
	for i, key := range keys {
		entry := timeSeries[key]
		entry.Date = key
		entry.Symbol = target.Metadata.Symbol
		dailyAdjustedTimeSeriesEntries[i] = entry
	}

	return &DailyAdjustedTimeSeries{Metadata: target.Metadata, TimeSeries: dailyAdjustedTimeSeriesEntries}, nil
}
