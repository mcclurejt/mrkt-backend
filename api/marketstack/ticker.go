package marketstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mcclurejt/mrkt-backend/api/common"
	db "github.com/mcclurejt/mrkt-backend/api/dynamodb"
)

type Exchange string

const (
	ExchangeNYSE    Exchange = "XNYS"
	ExchangeNasdaq  Exchange = "XNAS"
	DefaultExchange Exchange = ExchangeNYSE
	DefaultLimit             = 100
	MaxLimit                 = 1000
	DefaultOffset            = 0
)

const (
	TickerTableName = "Ticker"
)

type TickerResponse struct {
	Pagination TickerPagination `json:"pagination"`
	Data       []*TickerEntry   `json:"data"`
}

type TickerPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
	Total  int `json:"total"`
}

type TickerEntry struct {
	Symbol string `json:"symbol"`
}

type TickerService interface {
	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput

	Get(options *TickerOptions) ([]*TickerEntry, error)
	GetBatch(options *TickerOptions, ch chan<- common.ResultError)
}

type TickerOptions struct {
	Exchange Exchange
	Limit    int
	Offset   int
	Search   string
	sync.RWMutex
}

func DefaultTickerOptions() *TickerOptions {
	return &TickerOptions{Exchange: ExchangeNYSE, Limit: DefaultLimit, Offset: DefaultOffset}
}

func (o *TickerOptions) ToQueryString() string {
	o.Lock()
	defer o.Unlock()
	qs := fmt.Sprintf("&exchange=%s&limit=%d&offset=%d", o.Exchange, o.Limit, o.Offset)
	if o.Search != "" {
		qs += fmt.Sprintf("&search=%s", o.Search)
	}
	return qs
}

type tickerServicer struct {
	base baseClient
}

func newTickerService(base baseClient) TickerService {
	return tickerServicer{
		base: base,
	}
}

func (s tickerServicer) GetCreateTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Symbol"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Symbol"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(db.DefaultBillingMode),
		TableName:   aws.String(TickerTableName),
	}
}

func (s tickerServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(TickerTableName),
	}
}

func (s tickerServicer) Get(options *TickerOptions) ([]*TickerEntry, error) {
	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	ts, err := parseTickers(resp)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s tickerServicer) GetBatch(options *TickerOptions, ch chan<- common.ResultError) {
	tickers, err := s.Get(options)
	if err != nil {
		ch <- common.ResultError{Error: err}
	} else {
		ch <- common.ResultError{Result: tickers}
	}
}

func parseTickers(resp *http.Response) ([]*TickerEntry, error) {
	var target TickerResponse
	err := json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}

	return target.Data, nil
}
