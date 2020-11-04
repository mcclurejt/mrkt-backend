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
	Symbol    string `json:"symbol" an:"Symbol" at:"S" kt:"HASH"`
	Date      string `json:"date" an:"Date" at:"S" kt:"RANGE"`
	IsEnabled bool   `json:"isEnabled"`
}

func (s *IEXSymbolsServiceOp) Get(ctx context.Context) ([]IEXSymbol, error) {
	symbols := []IEXSymbol{}
	endpoint := fmt.Sprintf("/ref-data/iex/symbols")
	err := s.client.GetJSON(ctx, endpoint, &symbols)
	return symbols, err
}
