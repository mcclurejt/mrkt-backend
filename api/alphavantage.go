package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	DEFAULT_BASE_URL        = "https://www.alphavantage.co/query?"
	DEFAULT_TIMEOUT_SECONDS = 60
)

type AlphaVantageClient struct {
	base baseClient

	MonthlyAdjustedTimeSeriesService MonthlyAdjustedTimeSeriesService
}

func NewAlphaVantageClient(apiKey string) AlphaVantageClient {
	base := newAlphaVantageBaseClient(apiKey)
	return AlphaVantageClient{
		base:                             base,
		MonthlyAdjustedTimeSeriesService: newMonthlyAdjustedTimeSeriesService(base),
	}
}

type alphaVantageBaseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (a alphaVantageBaseClient) call(options requestOptions) (*http.Response, error) {
	url := DEFAULT_BASE_URL + fmt.Sprintf("apikey=%s", a.apiKey) + options.ToQueryString()
	return a.httpClient.Get(url)
}

func newAlphaVantageBaseClient(apiKey string) baseClient {
	return alphaVantageBaseClient{
		httpClient: &http.Client{
			Timeout: DEFAULT_TIMEOUT_SECONDS * time.Second,
		},
		apiKey: apiKey,
	}
}
