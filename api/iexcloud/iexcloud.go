package iexcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/mcclurejt/mrkt-backend/api/common"
)

type IexCloudRouteName string

const (
	IexCloudBaseURL = "https://cloud.iexapis.com/v1"
)

type baseClient struct {
	httpClient *http.Client
	url        string
	apiKey     string
}

type IexCloudClient struct {
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
}

type IexCloudError struct {
	StatusCode int
	Message    string
}

func (e IexCloudError) Error() string {
	return fmt.Sprintf("%d %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

func NewIexCloudClient(apiKey string, options ...func(*IexCloudClient)) *IexCloudClient {
	base := newIexCloudBaseClient(apiKey)
	c := &IexCloudClient{
		base: base,
	}

	for _, option := range options {
		option(c)
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

	return c
}

func newIexCloudBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: common.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
		url:    IexCloudBaseURL,
	}
}

func (c *IexCloudClient) GetJSON(ctx context.Context, endpoint string, v interface{}) error {
	addr, err := c.addToken(endpoint)
	if err != nil {
		return err
	}
	fmt.Printf("Address: %s\n", addr)
	data, err := c.getBytes(ctx, addr)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *IexCloudClient) GetJSONWithoutToken(ctx context.Context, endpoint string, v interface{}) error {
	addr := c.base.url + endpoint
	data, err := c.getBytes(ctx, addr)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *IexCloudClient) addToken(endpoint string) (string, error) {
	u, err := url.Parse(c.base.url + endpoint)
	if err != nil {
		return "", err
	}
	v := u.Query()
	v.Add("token", c.base.apiKey)
	u.RawQuery = v.Encode()
	return u.String(), nil
}

func (c *IexCloudClient) getBytes(ctx context.Context, addr string) ([]byte, error) {
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return []byte{}, err
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
		return []byte{}, IexCloudError{StatusCode: resp.StatusCode, Message: msg}
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *IexCloudClient) addOptions(s string, opt interface{}) (string, error) {
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
