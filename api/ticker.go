package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/mcclurejt/mrkt-backend/database/dynamodb"
)

const (
	EXCHANGE_NYSE   = "XNYS"
	EXCHANGE_NASDAQ = "XNAS"
)

const (
	TICKER_TABLE_NAME = "Ticker"
)

type Tickers struct {
	Data []string
}

type TickersResponse struct {
	Data []TickerEntry `json:data`
}

type TickerEntry struct {
	Symbol string `json:symbol`
}

type TickerService interface {
	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput

	Get(exchange string, limit int, offset int) ([]TickerEntry, error)
	Sync(exchange string, limit int, offset int, db db.Client) error
}

type tickerServiceOptions struct {
	Exchange string
	Limit    int
	Offset   int
}

func newTickerServiceOptions(exchange string, limit int, offset int) tickerServiceOptions {
	return tickerServiceOptions{Exchange: exchange, Limit: limit, Offset: offset}
}

func (o tickerServiceOptions) ToQueryString() string {
	return fmt.Sprintf("&exchange=%s&limit=%d&offset=%d", o.Exchange, o.Limit, o.Offset)
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
		TableName:   aws.String(TICKER_TABLE_NAME),
	}
}

func (s tickerServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(TICKER_TABLE_NAME),
	}
}

func (s tickerServicer) Get(exchange string, limit int, offset int) ([]TickerEntry, error) {
	options := newTickerServiceOptions(exchange, limit, offset)
	resp, err := s.base.call(options)
	if err != nil {
		return []TickerEntry{}, err
	}

	ts, err := parseTickers(resp)
	if err != nil {
		return []TickerEntry{}, err
	}
	return ts, nil
}

func (s tickerServicer) Sync(exchange string, limit int, offset int, db db.Client) error {
	//TODO
	return nil
}

func parseTickers(resp *http.Response) ([]TickerEntry, error) {
	target := &TickersResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return []TickerEntry{}, err
	}

	for i, t := range target.Data {
		target.Data[i].Symbol = strings.Replace(t.Symbol, ".", "-", -1)
	}
	return target.Data, nil
}
