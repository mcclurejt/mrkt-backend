package alphavantage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mcclurejt/mrkt-backend/api/base"
)

const (
	defaultBaseURL     = "https://www.alphavantage.co/query?"
	defaultRetryPeriod = 60
)

type AlphaVantageClient struct {
	base *baseClient

	MonthlyAdjustedTimeSeries MonthlyAdjustedTimeSeriesService
	CompanyOverview           CompanyOverviewService
	DailyAdjustedTimeSeries   DailyAdjustedTimeSeriesService
}

func NewAlphaVantageClient(apiKey string) AlphaVantageClient {
	base := newAlphaVantageBaseClient(apiKey)
	return AlphaVantageClient{
		base:                      base,
		MonthlyAdjustedTimeSeries: newMonthlyAdjustedTimeSeriesService(base),
		CompanyOverview:           newCompanyOverviewService(base),
		DailyAdjustedTimeSeries:   newDailyAdjustedTimeSeriesService(base),
	}
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (a *baseClient) call(options base.RequestOptions) (*http.Response, error) {
	url := defaultBaseURL + fmt.Sprintf("apikey=%s", a.apiKey) + options.ToQueryString()
	resp, err := a.httpClient.Get(url)
	return resp, err
}

func newAlphaVantageBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: base.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
