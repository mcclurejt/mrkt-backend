package marketstack

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/mcclurejt/mrkt-backend/api/common"
)

type MarketStackRouteName string

const (
	MarketStackBaseURL                      = "http://api.marketstack.com/v1/tickers"
	TickerRouteName    MarketStackRouteName = "ticker"
)

type MarketStackClient struct {
	base baseClient

	Ticker TickerService
}

func NewMarketStackClient(apiKey string) MarketStackClient {
	base := newMarketStackClient(apiKey)
	return MarketStackClient{
		base:   base,
		Ticker: newTickerService(base),
	}
}

func (m MarketStackClient) BatchCall(routeName MarketStackRouteName, target interface{}, options common.RequestOptions) error {
	switch routeName {
	case TickerRouteName:
		o, ok := options.(*TickerOptions)
		if !ok {
			return common.OptionParseError{DesiredType: reflect.TypeOf(&TickerOptions{})}
		}
		return m.GetBatchTickers(target, o)
	default:
		return &common.RouteNotRecognizedError{Route: string(routeName)}
	}
}

func (m MarketStackClient) GetBatchTickers(target interface{}, t *TickerOptions) error {
	ch := make(chan common.ResultError)
	totalLimit := t.Limit
	if t.Limit > MaxLimit {
		t = &TickerOptions{Limit: MaxLimit, Offset: t.Offset, Exchange: t.Exchange, Search: t.Search}
		count := 0
		for t.Limit > 0 {
			go func(ts *TickerOptions) {
				m.Ticker.GetBatch(ts, ch)
			}(t)
			// increment offset
			count += t.Limit
			t = &TickerOptions{Limit: t.Limit, Offset: t.Offset + t.Limit, Exchange: t.Exchange, Search: t.Search}
			if t.Offset+t.Limit > totalLimit {
				t = &TickerOptions{Limit: totalLimit - t.Offset, Offset: t.Offset, Exchange: t.Exchange, Search: t.Search}
			}
		}
	} else {
		go func(ts *TickerOptions) {
			m.Ticker.GetBatch(ts, ch)
		}(t)
	}
	entries := target.(*[]*TickerEntry)
	err := common.CollectResults(ch, int(totalLimit/MaxLimit+1), entries)
	if err != nil {
		return err
	}
	return nil
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (m baseClient) call(options common.RequestOptions) (*http.Response, error) {
	url := MarketStackBaseURL + fmt.Sprintf("?access_key=%s", m.apiKey) + options.ToQueryString()
	return m.httpClient.Get(url)
}

func newMarketStackClient(apiKey string) baseClient {
	return baseClient{
		httpClient: &http.Client{
			Timeout: common.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
