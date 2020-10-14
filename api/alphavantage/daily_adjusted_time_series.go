package alphavantage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mcclurejt/mrkt-backend/api/common"
	db "github.com/mcclurejt/mrkt-backend/api/dynamodb"
	"github.com/mcclurejt/mrkt-backend/database"
)

type OutputSize string

const (
	DailyAdjustedTimeSeriesFunction             = "TIME_SERIES_DAILY_ADJUSTED"
	DailyAdjustedTimeSeriesTableName            = "DailyAdjustedTimeSeries"
	OutputSizeCompact                OutputSize = "compact"
	OutputSizeFull                   OutputSize = "full"
	OutputSizeDefault                OutputSize = OutputSizeCompact
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
	Get(options *DailyAdjustedTimeSeriesOptions) (*DailyAdjustedTimeSeries, error)
	GetBatch(options *DailyAdjustedTimeSeriesOptions, ch chan<- common.ResultError)
	Sync(options *DailyAdjustedTimeSeriesOptions, db database.SQLClient) error

	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput
}

type DailyAdjustedTimeSeriesOptions struct {
	Symbol     string
	OutputSize OutputSize
	sync.RWMutex
}

func (o *DailyAdjustedTimeSeriesOptions) ToQueryString() string {
	o.Lock()
	defer o.Unlock()
	return fmt.Sprintf("&function=%s&symbol=%s&outputsize=%s", DailyAdjustedTimeSeriesFunction, o.Symbol, o.OutputSize)
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
		TableName:   aws.String(DailyAdjustedTimeSeriesTableName),
	}
}

func (s dailyAdjustedTimeSeriesServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(DailyAdjustedTimeSeriesTableName),
	}
}

func (s dailyAdjustedTimeSeriesServicer) Get(options *DailyAdjustedTimeSeriesOptions) (*DailyAdjustedTimeSeries, error) {
	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	ts, err := parseDailyAdjustedTimeSeries(resp)
	if err != nil {
		_, ok := err.(*AlphaVantageRateExceededError)
		if ok {
			time.Sleep(defaultRetryPeriod * time.Second)
			return s.Get(options)
		}
		return nil, err
	}

	return ts, nil
}

func (s dailyAdjustedTimeSeriesServicer) GetBatch(options *DailyAdjustedTimeSeriesOptions, ch chan<- common.ResultError) {
	ts, err := s.Get(options)
	if err != nil {
		ch <- common.ResultError{Error: err}
	} else {
		ch <- common.ResultError{Result: ts}
	}
}

func (s dailyAdjustedTimeSeriesServicer) Sync(options *DailyAdjustedTimeSeriesOptions, db database.SQLClient) error {
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
