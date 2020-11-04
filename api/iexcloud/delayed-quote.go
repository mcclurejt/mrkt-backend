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
	Symbol           string  `json:"symbol" dynamodbav:"Symbol" attributetype:"S" keytype:"HASH"`
	DelayedPriceTime string  `dynamodbav:"Time" attributetype:"S" keytype:"RANGE"`
	delayedPriceTime int64   `json:"delayedPriceTime"`
	DelayedPrice     float64 `json:"delayedPrice" dynamodbav:"PriceTime"`
	DelayedSize      int     `json:"delayedSize" dynamodbav:"Size"`
	High             float64 `json:"High" dynamodbav:"High"`
	Low              float64 `json:"Low" dynamodbav:"Low"`
	TotalVolume      int     `json:"totalVolume" dynamodbav:"Volume"`
	ProcessedTime    int     `json:"processedTime" dynamodbav:"ProcessedTime"`
}

func (s *DelayedQuoteServiceOp) Get(ctx context.Context, symbol string) (*DelayedQuote, error) {
	delayedQuote := &DelayedQuote{}
	endpoint := fmt.Sprintf("/stock/%s/delayed-quote", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, delayedQuote)
	delayedQuote.DelayedPriceTime = TimeToTimestamp(delayedQuote.delayedPriceTime)
	return delayedQuote, err
}
