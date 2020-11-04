package iexcloud

import (
	"context"
	"fmt"
)

type IexSymbolsService interface {
	Get(ctx context.Context) ([]IexSymbol, error)
}

type IexSymbolsServiceOp struct {
	client *IexCloudClient
}

var _ IexSymbolsService = &IexSymbolsServiceOp{}

type IexSymbol struct {
	Symbol    string `json:"symbol"`
	Date      string `json:"date"`
	IsEnabled bool   `json:"isEnabled"`
}

func (s *IexSymbolsServiceOp) Get(ctx context.Context) ([]IexSymbol, error) {
	symbols := []IexSymbol{}
	endpoint := fmt.Sprintf("/ref-data/iex/symbols")
	err := s.client.GetJSON(ctx, endpoint, &symbols)
	return symbols, err
}
