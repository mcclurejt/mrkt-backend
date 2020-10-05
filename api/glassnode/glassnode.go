package glassnode

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mcclurejt/mrkt-backend/api/base"
)

const (
	GLASSNODE_BASE_URL = "https://api.glassnode.com/v1/metrics/"
)

type GlassNodeClient struct {
	base *baseClient

	Coin                    CoinService
	NetUnrealizedProfitLoss NetUnrealizedProfitLossService
}

func NewGlassNodeClient(apiKey string) GlassNodeClient {
	base := newGlassNodeBaseClient(apiKey)
	return GlassNodeClient{
		base:                    base,
		Coin:                    newCoinService(base),
		NetUnrealizedProfitLoss: newNetUnrealizedProfitLossService(base),
	}
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (g baseClient) call(options base.RequestOptions) (*http.Response, error) {
	url := GLASSNODE_BASE_URL + options.ToQueryString() + fmt.Sprintf("&api_key=%s", g.apiKey)
	resp, err := g.httpClient.Get(url)
	return resp, err
}

func newGlassNodeBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: base.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
