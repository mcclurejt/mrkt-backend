package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type LargestTradesService interface {
	Get(ctx context.Context, symbol string) (*[]LargestTrade, error)
}

type LargestTradesServiceOp struct {
	client *IexCloudClient
}

var _ LargestTradesService = &LargestTradesServiceOp{}

type LargestTrade struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Time      int     `json:"time"`
	TimeLabel string  `json:"timeLabel"`
	Venue     string  `json:"venue"`
	VenueName string  `json:"venueName"`
}

func (s *LargestTradesServiceOp) Get(ctx context.Context, symbol string) (*[]LargestTrade, error) {
	lt := new([]LargestTrade)
	endpoint := fmt.Sprintf("/stock/%s/largest-trades", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &lt)
	return lt, err
}
