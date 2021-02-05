package iex

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
)

const (
	DefaultClientTimeout = 10
	DefaultBaseURL       = "https://cloud.iexapis.com/v1"
)

type Client interface {
	ChartService
}

type RealClient struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string

	Log *logrus.Logger
}

var _ Client = &RealClient{}

func New(apiKey string, options ...ClientOption) Client {
	base := &RealClient{
		HTTPClient: http.DefaultClient,
		BaseURL:    DefaultBaseURL,
		APIKey:     apiKey,
		Log:        logrus.New(),
	}
	for _, opt := range options {
		base = opt(base)
	}
	return base
}

func (c *RealClient) get(ctx context.Context, path string, v interface{}) error {
	return c.getWithParams(ctx, path, v, nil)
}

func (c *RealClient) getWithParams(ctx context.Context, path string, v interface{}, queryParamStruct interface{}) error {
	// build the url
	url, err := encodeQueryParams(c.BaseURL+path, queryParamStruct)
	if err != nil {
		return err
	}
	queryParams := url.Query()
	// add api token
	queryParams["token"] = []string{c.APIKey}
	url.RawQuery = queryParams.Encode()
	// execute the request
	c.Log.WithFields(logrus.Fields{"url": url.String()}).Debug("Executing GET request")
	resp, err := c.HTTPClient.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// check the response status
	bytes, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if err != nil {
			return err
		}
		return errors.New(string(bytes))
	}
	// decode the response body
	c.Log.Debugf("Response Body: %s", bytes)
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

func encodeQueryParams(base string, opt interface{}) (*url.URL, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		u, err := url.Parse(base)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	origURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	origValues := origURL.Query()
	newValues, err := query.Values(opt)
	if err != nil {
		return nil, err
	}
	for k, v := range newValues {
		origValues[k] = v
	}
	origURL.RawQuery = origValues.Encode()
	return origURL, nil
}
