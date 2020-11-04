package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type IntradayPricesService interface {
	Get(ctx context.Context, symbol string) ([]IntradayPrice, error)
	GetWithOptions(ctx context.Context, symbol string, options *IntradayOptions) ([]IntradayPrice, error)
}

type IntradayPricesServiceOp struct {
	client *IEXCloudClient
}

var _ IntradayPricesService = &IntradayPricesServiceOp{}

type IntradayPrice struct {
	Date                 string  `json:"date"`
	Minute               string  `json:"minute"`
	Label                string  `json:"label"`
	MarketOpen           float64 `json:"marketOpen"`
	MarketClose          float64 `json:"marketClose"`
	MarketHigh           float64 `json:"marketHigh"`
	MarketLow            float64 `json:"marketLow"`
	MarketAverage        float64 `json:"marketAverage"`
	MarketVolume         int     `json:"marketVolume"`
	MarketNotional       float64 `json:"marketNotional"`
	MarketNumTrades      int     `json:"marketNumberOfTrades"`
	MarketChangeOverTime float64 `json:"marketChangeOverTime"`
	High                 float64 `json:"High"`
	Low                  float64 `json:"Low"`
	Open                 float64 `json:"Open"`
	Close                float64 `json:"Close"`
	Average              float64 `json:"average"`
	Volume               int     `json:"volume"`
	Notional             float64 `json:"notional"`
	NumTrades            int     `json:"numberOfTrades"`
	ChangeOverTime       float64 `json:"changeOverTime"`
}

type IntradayOptions struct {
	ChartSimplify   bool   `url:"chartSimplify,omitempty"`
	ChartInterval   int    `url:"chartInterval,omitempty"`
	ChangeFromClose bool   `url:"changeFromClose,omitempty"`
	ChartLast       int    `url:"chartLast,omitempty"`
	ExactDate       string `url:"exactDate,omitempty"` // Formatted as YYYYMMDD
}

func (s *IntradayPricesServiceOp) Get(ctx context.Context, symbol string) ([]IntradayPrice, error) {
	intradayPrices := []IntradayPrice{}
	endpoint := fmt.Sprintf("/stock/%s/intraday-prices", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &intradayPrices)
	return intradayPrices, err
}

func (s *IntradayPricesServiceOp) GetWithOptions(ctx context.Context, symbol string, options *IntradayOptions) ([]IntradayPrice, error) {
	intradayPrices := []IntradayPrice{}
	endpoint := fmt.Sprintf("/stock/%s/intraday-prices", url.PathEscape(symbol))
	endpoint, err := s.client.addOptions(endpoint, options)
	err = s.client.GetJSON(ctx, endpoint, &intradayPrices)
	return intradayPrices, err
}
