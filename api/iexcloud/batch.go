package iexcloud

import (
	"context"
	"fmt"
)

// Query Types

type QueryType string

const (
	QueryTypeBook           QueryType = "book"
	QueryTypeDelayedQuote   QueryType = "delayed-quote"
	QueryTypeIntradayPrices QueryType = "intraday-prices"
	QueryTypeLargestTrades  QueryType = "largest-trades"
	// OHLC = "ohlc" // open-high-low-close
	// PREVIOUS = "previous"
	// QUOTE = "quote"

	QueryTypeCompany             QueryType = "company"
	QueryTypeInsiderRoster       QueryType = "insider-roster"
	QueryTypeInsiderSummary      QueryType = "insider-summary"
	QueryTypeInsiderTransactions QueryType = "insider-transactions"
	QueryTypePeers               QueryType = "peers"
)

type BatchService interface {
	GetMarketBatch(ctx context.Context, symbols []string, types []QueryType) (map[string]Batch, error)
	GetSymbolBatch(ctx context.Context, symbol string, types []QueryType) (*Batch, error)
}

type BatchServiceOp struct {
	client *IEXCloudClient
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

func (s *BatchServiceOp) GetMarketBatch(ctx context.Context, symbols []string, types []QueryType) (map[string]Batch, error) {
	batch := map[string]Batch{}
	endpoint := fmt.Sprintf("/stock/market/batch/")
	options := &BatchOptions{
		Symbols: SliceToString(symbols, StrToPtr(",")),
		Types:   SliceToString(types, StrToPtr(",")),
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
		Types: SliceToString(types, StrToPtr(",")),
	}
	endpoint, err := s.client.addOptions(endpoint, options)
	if err != nil {
		return nil, err
	}
	err = s.client.GetJSON(ctx, endpoint, batch)
	return batch, err
}
