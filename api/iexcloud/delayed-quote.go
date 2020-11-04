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
	Symbol           string  `json:"symbol" attributetype:"S" keytype:"HASH"`
	DelayedPriceTime string  `dynamodbav:"time" attributetype:"S" keytype:"RANGE"`
	delayedPriceTime int64   `json:"delayedPriceTime"`
	DelayedPrice     float64 `json:"delayedPrice" dynamodbav:"price"`
	DelayedSize      int     `json:"delayedSize" dynamodbav:"size"`
	High             float64 `json:"High" dynamodbav:"high"`
	Low              float64 `json:"Low" dynamodbav:"low"`
	TotalVolume      int     `json:"totalVolume" dynamodbav:"volume"`
	ProcessedTime    int     `json:"processedTime" dynamodbav:"processedTime"`
}

func (s *DelayedQuoteServiceOp) Get(ctx context.Context, symbol string) (*DelayedQuote, error) {
	delayedQuote := &DelayedQuote{}
	endpoint := fmt.Sprintf("/stock/%s/delayed-quote", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, delayedQuote)
	delayedQuote.DelayedPriceTime = TimeToTimestamp(delayedQuote.delayedPriceTime)
	return delayedQuote, err
}
