package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	MARKETSTACK_BASE_URL = "http://api.marketstack.com/v1/tickers"
)

type MarketStackClient struct {
	base baseClient

	TickerService TickerService
}

func NewMarketStackClient(apiKey string) MarketStackClient {
	base := newMarketStackClient(apiKey)
	return MarketStackClient{
		base:          base,
		TickerService: newTickerService(base),
	}
}

type marketStackBaseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (m marketStackBaseClient) call(options requestOptions) (*http.Response, error) {
	url := MARKETSTACK_BASE_URL + fmt.Sprintf("?access_key=%s", m.apiKey) + options.ToQueryString()
	return m.httpClient.Get(url)
}

func newMarketStackClient(apiKey string) baseClient {
	return marketStackBaseClient{
		httpClient: &http.Client{
			Timeout: DEFAULT_TIMEOUT_SECONDS * time.Second,
		},
		apiKey: apiKey,
	}
}
