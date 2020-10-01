package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	TICKER_TABLE_NAME = "Ticker"
)

var (
	TICKER_HEADERS = []string{
		"id",
		"name",
	}
	TICKER_COLUMNS = []string{
		"int INT NOT NULL AUTO_INCREMENT PRIMARY KEY",
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
	Get(exchange string, limit int, offset int) (Tickers, error)
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
