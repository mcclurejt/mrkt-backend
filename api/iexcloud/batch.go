package iexcloud

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// Query Types

type QueryType string

const (
	BOOK            QueryType = "book"
	DELAYED_QUOTE   QueryType = "delayed-quote"
	INTRADAY_PRICES QueryType = "intraday-prices"
	LARGEST_TRADES  QueryType = "largest-trades"
	// OHLC = "ohlc" // open-high-low-close
	// PREVIOUS = "previous"
	// QUOTE = "quote"

	COMPANY              QueryType = "company"
	INSIDER_ROSTER       QueryType = "insider-roster"
	INSIDER_SUMMARY      QueryType = "insider-summary"
	INSIDER_TRANSACTIONS QueryType = "insider-transactions"
	PEERS                QueryType = "peers"
)

type BatchService interface {
	GetMarketBatch(ctx context.Context, symbols []string, types []QueryType) (*map[string]Batch, error)
	GetSymbolBatch(ctx context.Context, symbol string, types []QueryType) (*Batch, error)
}

type BatchServiceOp struct {
	client *IexCloudClient
}

var _ BatchService = &BatchServiceOp{}

type Batch struct {
	Book                Book                 `json:"book,omitEmpty"`
	Quote               Quote                `json:"quote,omitEmpty"`
	InsiderRoster       []InsiderRoster      `json:"insider-roster,omitEmpty"`
	InsiderSummary      []InsiderSummary     `json:"insider-summary,omitEmpty"`
	InsiderTransactions []InsiderTransaction `json:"insider-transactions,omitEmpty"`
	Company             Company              `json:"company,omitEmpty"`
}

type BatchOptions struct {
	Symbols string `url:"symbols,omitEmpty"`
	Types   string `url:"types,omitEmpty"`
}

type SymbolBatchOptions struct {
	Types string `url:"types,omitEmpty"`
}

func (s *BatchServiceOp) GetMarketBatch(ctx context.Context, symbols []string, types []QueryType) (*map[string]Batch, error) {
	batch := new(map[string]Batch)
	endpoint := fmt.Sprintf("/stock/market/batch/")
	options := &BatchOptions{
		Symbols: ToURLString(symbols),
		Types:   ToURLString(types),
	}
	endpoint, err := s.client.addOptions(endpoint, options)
	if err != nil {
		return nil, err
	}
	err = s.client.GetJSON(ctx, endpoint, &batch)
	return batch, err
}

func (s *BatchServiceOp) GetSymbolBatch(ctx context.Context, symbol string, types []QueryType) (*Batch, error) {
	batch := new(Batch)
	endpoint := fmt.Sprintf("/stock/%s/batch/", symbol)
	options := &SymbolBatchOptions{
		Types: ToURLString(types),
	}
	endpoint, err := s.client.addOptions(endpoint, options)
	if err != nil {
		return nil, err
	}
	err = s.client.GetJSON(ctx, endpoint, &batch)
	return batch, err
}

// ToURLString - takes a slice of string-like objects and converts them to a string containing the comma-separated items
func ToURLString(arr interface{}) string {
	t := reflect.TypeOf(arr)
	if t.Kind() != reflect.Slice {
		panic(arr)
	}
	v := reflect.ValueOf(arr)
	l := v.Len()
	stringArr := make([]string, l)
	for i := 0; i < l; i++ {
		entry := v.Index(i)
		stringArr[i] = entry.String()
	}
	return strings.Join(stringArr, ",")
}
