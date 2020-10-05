package marketstack

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mcclurejt/mrkt-backend/api/base"
)

const (
	MARKETSTACK_BASE_URL = "http://api.marketstack.com/v1/tickers"
)

type MarketStackClient struct {
	base *baseClient

	Ticker TickerService
}

func NewMarketStackClient(apiKey string) MarketStackClient {
	base := newMarketStackClient(apiKey)
	return MarketStackClient{
		base:   base,
		Ticker: newTickerService(base),
	}
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (m baseClient) call(options base.RequestOptions) (*http.Response, error) {
	url := MARKETSTACK_BASE_URL + fmt.Sprintf("?access_key=%s", m.apiKey) + options.ToQueryString()
	return m.httpClient.Get(url)
}

func newMarketStackClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: base.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
