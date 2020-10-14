package alphavantage

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/mcclurejt/mrkt-backend/api/common"
)

type AlphaVantageRouteName string

const (
	defaultBaseURL                                   = "https://www.alphavantage.co/query?"
	defaultRetryPeriod                               = 60
	DailyTimeSeriesRouteName   AlphaVantageRouteName = "dailyadjustedtimeseries"
	MonthlyTimeSeriesRouteName AlphaVantageRouteName = "monthlyadjustedtimeseries"
	CompanyOverviewRouteName   AlphaVantageRouteName = "companyoverview"
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

func (c *AlphaVantageClient) BatchCall(routeName AlphaVantageRouteName, assets []string, target interface{}, options common.RequestOptions) error {
	switch routeName {
	case DailyTimeSeriesRouteName:
		o, ok := options.(*DailyAdjustedTimeSeriesOptions)
		if !ok {
			return common.OptionParseError{DesiredType: reflect.TypeOf(&DailyAdjustedTimeSeriesOptions{})}
		}
		return c.GetBatchDailyTimeSeries(assets, target, o)
	case MonthlyTimeSeriesRouteName:
		_, ok := options.(*MonthlyAdjustedTimeSeriesOptions)
		if !ok {
			return common.OptionParseError{DesiredType: reflect.TypeOf(&MonthlyAdjustedTimeSeriesOptions{})}
		}
	case CompanyOverviewRouteName:
		_, ok := options.(*CompanyOverviewOptions)
		if !ok {
			return common.OptionParseError{DesiredType: reflect.TypeOf(&CompanyOverviewOptions{})}
		}
	default:
		return &common.RouteNotRecognizedError{Route: string(routeName)}
	}
	return nil
}

func (c *AlphaVantageClient) GetBatchDailyTimeSeries(assets []string, target interface{}, options *DailyAdjustedTimeSeriesOptions) error {
	ch := make(chan common.ResultError)
	for _, asset := range assets {
		go func(o *DailyAdjustedTimeSeriesOptions) {
			c.DailyAdjustedTimeSeries.GetBatch(o, ch)
		}(&DailyAdjustedTimeSeriesOptions{Symbol: asset, OutputSize: options.OutputSize})
	}
	entries := target.(*[]*DailyAdjustedTimeSeriesEntry)
	err := common.CollectResults(ch, len(assets), entries)
	if err != nil {
		return err
	}
	return nil
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (a *baseClient) call(options common.RequestOptions) (*http.Response, error) {
	url := defaultBaseURL + fmt.Sprintf("apikey=%s", a.apiKey) + options.ToQueryString()
	resp, err := a.httpClient.Get(url)
	return resp, err
}

func newAlphaVantageBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: common.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
