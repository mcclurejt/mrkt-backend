package iexcloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
	"golang.org/x/time/rate"
)

type IEXCloudRouteName string

const (
	DefaultTimeout  = 10
	IEXCloudBaseURL = "https://cloud.iexapis.com/v1"
)

var ErrorEmptyResponse = errors.New("Response Body Was Empty")

type baseClient struct {
	httpClient  *http.Client
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
}

type IEXCloudClient struct {
	base *baseClient

	// services
	Status              StatusService
	Book                BookService
	DelayedQuote        DelayedQuoteService
	IntradayPrices      IntradayPricesService
	LargestTrades       LargestTradesService
	Company             CompanyService
	InsiderRoster       InsiderRosterService
	InsiderSummary      InsiderSummaryService
	InsiderTransactions InsiderTransactionsService
	Peers               PeersService
	Batch               BatchService
	SectorPerformance   SectorPerformanceService
	Options             OptionsService
	IexSymbols          IEXSymbolsService
	Chart               ChartService
}

type IEXCloudError struct {
	StatusCode int
	Message    string
}

func (e IEXCloudError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

func NewIEXCloudClient(apiKey string, options ...func(*IEXCloudClient)) *IEXCloudClient {
	base := newIEXCloudBaseClient(apiKey)
	c := &IEXCloudClient{
		base: base,
	}
	// services
	c.Status = &StatusServiceOp{client: c}
	c.Book = &BookServiceOp{client: c}
	c.DelayedQuote = &DelayedQuoteServiceOp{client: c}
	c.IntradayPrices = &IntradayPricesServiceOp{client: c}
	c.LargestTrades = &LargestTradesServiceOp{client: c}
	c.Company = &CompanyServiceOp{client: c}
	c.InsiderRoster = &InsiderRosterServiceOp{client: c}
	c.InsiderSummary = &InsiderSummaryServiceOp{client: c}
	c.InsiderTransactions = &InsiderTransactionsServiceOp{client: c}
	c.Peers = &PeersServiceOp{client: c}
	c.Batch = &BatchServiceOp{client: c}
	c.SectorPerformance = &SectorPerformanceServiceOp{client: c}
	c.Options = &OptionsServiceOp{client: c}
	c.IexSymbols = &IEXSymbolsServiceOp{client: c}
	c.Chart = &ChartServiceOp{client: c}

	for _, option := range options {
		option(c)
	}

	return c
}

func newIEXCloudBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: DefaultTimeout * time.Second,
		},
		apiKey:      apiKey,
		url:         IEXCloudBaseURL,
		rateLimiter: rate.NewLimiter(10, 1),
	}
}

func (c *IEXCloudClient) GetJSON(ctx context.Context, endpoint string, v interface{}) error {
	addr, err := c.addToken(endpoint)
	if err != nil {
		return err
	}
	data, err := c.getBytes(ctx, addr)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *IEXCloudClient) GetJSONWithoutToken(ctx context.Context, endpoint string, v interface{}) error {
	addr := c.base.url + endpoint
	data, err := c.getBytes(ctx, addr)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *IEXCloudClient) addToken(endpoint string) (string, error) {
	u, err := url.Parse(c.base.url + endpoint)
	if err != nil {
		return "", err
	}
	v := u.Query()
	v.Add("token", c.base.apiKey)
	u.RawQuery = v.Encode()
	return u.String(), nil
}

func (c *IEXCloudClient) getBytes(ctx context.Context, addr string) ([]byte, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return []byte{}, err
	}
	err = c.base.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.base.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		msg := ""
		if err == nil {
			msg = string(b)
		}
		return []byte{}, IEXCloudError{StatusCode: resp.StatusCode, Message: msg}
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *IEXCloudClient) addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}
