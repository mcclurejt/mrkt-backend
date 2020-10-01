package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	GLASSNODE_BASE_URL = "https://api.glassnode.com/v1/metrics/"
)

type GlassNodeClient struct {
	base baseClient

	NetUnrealizedProfitLossService NetUnrealizedProfitLossService
}

func NewGlassNodeClient(apiKey string) GlassNodeClient {
	base := newGlassNodeBaseClient(apiKey)
	return GlassNodeClient{
		base:                           base,
		NetUnrealizedProfitLossService: newNetUnrealizedProfitLossService(base),
	}
}

type glassNodeBaseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (g glassNodeBaseClient) call(options requestOptions) (*http.Response, error) {
	url := GLASSNODE_BASE_URL + options.ToQueryString() + fmt.Sprintf("&apikey=%s", g.apiKey)
	resp, err := g.httpClient.Get(url)
	return resp, err
}

func newGlassNodeBaseClient(apiKey string) baseClient {
	return glassNodeBaseClient{
		httpClient: &http.Client{
			Timeout: DEFAULT_TIMEOUT_SECONDS * time.Second,
		},
		apiKey: apiKey,
	}
}
