package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	DEFAULT_BASE_URL             = "https://www.alphavantage.co/query?"
	DEFAULT_TIMEOUT_SECONDS      = 60
	DEFAULT_RETRY_PERIOD_SECONDS = 60
)

type AlphaVantageClient struct {
	base baseClient

	MonthlyAdjustedTimeSeriesService MonthlyAdjustedTimeSeriesService
	CompanyOverviewService           CompanyOverviewService
}

func NewAlphaVantageClient(apiKey string) AlphaVantageClient {
	base := newAlphaVantageBaseClient(apiKey)
	return AlphaVantageClient{
		base:                             base,
		MonthlyAdjustedTimeSeriesService: newMonthlyAdjustedTimeSeriesService(base),
		CompanyOverviewService:           newCompanyOverviewService(base),
	}
}

type alphaVantageBaseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (a alphaVantageBaseClient) call(options requestOptions) (*http.Response, error) {
	url := DEFAULT_BASE_URL + fmt.Sprintf("apikey=%s", a.apiKey) + options.ToQueryString()
	resp, err := a.httpClient.Get(url)
	return resp, err
}

func newAlphaVantageBaseClient(apiKey string) baseClient {
	return alphaVantageBaseClient{
		httpClient: &http.Client{
			Timeout: DEFAULT_TIMEOUT_SECONDS * time.Second,
		},
		apiKey: apiKey,
	}
}
