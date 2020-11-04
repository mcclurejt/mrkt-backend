package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type DelayedQuoteService interface {
	Get(ctx context.Context, symbol string) (*DelayedQuote, error)
}

type DelayedQuoteServiceOp struct {
	client *IEXCloudClient
}

var _ DelayedQuoteService = &DelayedQuoteServiceOp{}

type DelayedQuote struct {
	Symbol           string  `json:"symbol"`
	DelayedPrice     float64 `json:"delayedPrice"`
	DelayedSize      int     `json:"delayedSize"`
	DelayedPriceTime int     `json:"delayedPriceTime"`
	High             float64 `json:"High"`
	Low              float64 `json:"Low"`
	TotalVolume      int     `json:"totalVolume"`
	ProcessedTime    int     `json:"processedTime"`
}

func (s *DelayedQuoteServiceOp) Get(ctx context.Context, symbol string) (*DelayedQuote, error) {
	delayedQuote := new(DelayedQuote)
	endpoint := fmt.Sprintf("/stock/%s/delayed-quote", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, delayedQuote)
	return delayedQuote, err
}
