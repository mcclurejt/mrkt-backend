package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	EXCHANGE_NYSE   = "XNYS"
	EXCHANGE_NASDAQ = "XNAS"
)

const (
	TICKER_TABLE_NAME = "Ticker"
)

var (
	TICKER_HEADERS = []string{
		"id",
		"name",
	}
	TICKER_INSERTION_HEADERS = []string{
		"name",
	}
	TICKER_COLUMNS = []string{
		"id INT NOT NULL AUTO_INCREMENT PRIMARY KEY",
		"name VARCHAR (8) NOT NULL UNIQUE",
	}
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
	GetTableName() string
	GetTableColumns() []string
	Get(exchange string, limit int, offset int) (Tickers, error)
	Insert(t Tickers, db database.SQLClient) error
	Sync(exchange string, limit int, offset int, db database.SQLClient) error
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

func (s tickerServicer) GetTableName() string {
	return TICKER_TABLE_NAME
}

func (s tickerServicer) GetTableColumns() []string {
	return TICKER_COLUMNS
}

func (s tickerServicer) Get(exchange string, limit int, offset int) (Tickers, error) {
	options := newTickerServiceOptions(exchange, limit, offset)
	resp, err := s.base.call(options)
	if err != nil {
		return Tickers{}, err
	}

	ts, err := parseTickers(resp)
	if err != nil {
		return Tickers{}, err
	}
	return ts, nil
}

func (s tickerServicer) Insert(t Tickers, db database.SQLClient) error {
	values := make([]interface{}, len(t.Data))
	for i, v := range t.Data {
		values[i] = v
	}
	return db.Insert(s.GetTableName(), TICKER_INSERTION_HEADERS, values)
}

func (s tickerServicer) Sync(exchange string, limit int, offset int, db database.SQLClient) error {
	//TODO
	return nil
}

func parseTickers(resp *http.Response) (Tickers, error) {
	target := &TickersResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return Tickers{}, err
	}

	ts := make([]string, len(target.Data))
	for i, t := range target.Data {
		ts[i] = t.Symbol
	}
	return Tickers{Data: ts}, nil
}
