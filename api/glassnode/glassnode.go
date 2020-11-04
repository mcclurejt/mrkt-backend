package glassnode

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type GlassNodeRouteName string

const (
	GlassNodeBaseURL                    = "https://api.glassnode.com"
	NuplRouteName    GlassNodeRouteName = "nupl"
)

type GlassNodeClient struct {
	base                    *baseClient
	Assets                  AssetService
	NetUnrealizedProfitLoss NetUnrealizedProfitLossService
}

func NewGlassNodeClient(apiKey string) *GlassNodeClient {
	base := newGlassNodeBaseClient(apiKey)
	return &GlassNodeClient{
		base:                    base,
		Assets:                  newAssetService(base),
		NetUnrealizedProfitLoss: newNetUnrealizedProfitLossService(base),
	}
}

func (c *GlassNodeClient) BatchCall(routeName GlassNodeRouteName, assets []string, target interface{}, options common.RequestOptions) error {
	switch routeName {
	case NuplRouteName:
		o, ok := options.(*NetUnrealizedProfitLossOptions)
		if !ok {
			return OptionParseError{DesiredType: reflect.TypeOf(&NetUnrealizedProfitLossOptions{})}
		}
		return c.GetBatchNupl(assets, target, o)

	default:
		return &RouteNotRecognizedError{Route: string(routeName)}
	}
}

func (c *GlassNodeClient) GetBatchNupl(assets []string, target interface{}, options *NetUnrealizedProfitLossOptions) error {
	ch := make(chan common.ResultError)
	for _, asset := range assets {
		options = &NetUnrealizedProfitLossOptions{
			Asset:    asset,
			Interval: options.Interval,
			Since:    options.Since,
			Until:    options.Until,
		}
		go func(o *NetUnrealizedProfitLossOptions) {
			c.NetUnrealizedProfitLoss.GetBatch(o, ch)
		}(&NetUnrealizedProfitLossOptions{
			Asset:    asset,
			Interval: options.Interval,
			Since:    options.Since,
			Until:    options.Until})
	}
	t := target.([]*NetUnrealizedProfitLossEntry)
	err := common.CollectResults(ch, len(assets), &t)
	if err != nil {
		return err
	}
	return nil
}

type baseClient struct {
	httpClient *http.Client
	apiKey     string
}

func (g baseClient) call(options common.RequestOptions) (*http.Response, error) {
	url := GlassNodeBaseURL + options.ToQueryString() + fmt.Sprintf("&api_key=%s", g.apiKey)
	return g.httpClient.Get(url)
}

func newGlassNodeBaseClient(apiKey string) *baseClient {
	return &baseClient{
		httpClient: &http.Client{
			Timeout: common.DefaultTimeout * time.Second,
		},
		apiKey: apiKey,
	}
}
