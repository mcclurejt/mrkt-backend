package iexcloud

import (
	"context"
	"fmt"
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
	GetMarketBatch(context.Context, []string, []string) (*map[string]Batch, error)
	GetSymbolBatch(context.Context, string, []string) (*Batch, error)
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

func (s *BatchServiceOp) GetMarketBatch(ctx context.Context, symbols []string, types []string) (*map[string]Batch, error) {
	batch := new(map[string]Batch)
	endpoint := fmt.Sprintf("/stock/market/batch/")
	options := &BatchOptions{
		Symbols: toQueryString(symbols),
		Types:   toQueryString(types),
	}
	endpoint, err := s.client.addOptions(endpoint, options)
	if err != nil {
		return nil, err
	}
	err = s.client.GetJSON(ctx, endpoint, &batch)
	return batch, err
}

func (s *BatchServiceOp) GetSymbolBatch(ctx context.Context, symbol string, types []string) (*Batch, error) {
	batch := new(Batch)
	endpoint := fmt.Sprintf("/stock/%s/batch/", symbol)
	options := &SymbolBatchOptions{
		Types: toQueryString(types),
	}
	endpoint, err := s.client.addOptions(endpoint, options)
	if err != nil {
		return nil, err
	}
	err = s.client.GetJSON(ctx, endpoint, &batch)
	return batch, err
}

func toQueryString(arr []string) string {
	return strings.Join(arr, ",")
}
