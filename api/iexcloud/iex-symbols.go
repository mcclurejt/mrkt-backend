package iexcloud

import (
	"context"
	"fmt"
)

type IEXSymbolsService interface {
	Get(ctx context.Context) ([]IEXSymbol, error)
}

type IEXSymbolsServiceOp struct {
	client *IEXCloudClient
}

var _ IEXSymbolsService = &IEXSymbolsServiceOp{}

type IEXSymbol struct {
	Symbol    string `json:"symbol" at:"S" kt:"HASH"`
	Date      string `at:"S" kt:"HASH"`
	date      string `json:"date"`
	isEnabled bool   `json:"isEnabled"`
}

func (s *IEXSymbolsServiceOp) Get(ctx context.Context) ([]IEXSymbol, error) {
	symbols := []IEXSymbol{}
	endpoint := fmt.Sprintf("/ref-data/iex/symbols")
	err := s.client.GetJSON(ctx, endpoint, &symbols)
	for i, symbol := range symbols {
		ts, err := DateToTimestamp(symbol.Date)
		if err == nil {
			symbols[i].Date = ts
		} else {
			symbols[i].Date = symbol.date
		}
	}
	return symbols, err
}
