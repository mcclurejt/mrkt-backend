package glassnode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	BTC = "BTC"
	ETH = "ETH"
)

const (
	assetAPIRoute        = "/v2/metrics/endpoints"
	assetTableName       = "Assets"
	defaultEndpointRoute = "/v1/metrics/indicators/net_unrealized_profit_loss"
)

type AssetResponse struct {
	Assets []*Asset
}

type Asset struct {
	Symbol          string   `json:"symbol"`
	Name            string   `json:"name"`
	IsERC20         bool     `json:"isERC20"`
	IsStablecoin    bool     `json:"isStablecoin"`
	IsExchangeToken bool     `json:"isExchangeToken"`
	Tags            []string `json:"tags"`
}

type AssetService interface {
	Get(options *AssetOptions) ([]*Asset, error)
	Sync(db database.SQLClient) error
}

type AssetOptions struct {
	Route string
}

func DefaultAssetOptions() *AssetOptions {
	return &AssetOptions{Route: defaultEndpointRoute}
}

func (o *AssetOptions) ToQueryString() string {
	return fmt.Sprintf("%s", assetAPIRoute)
}

type assetServicer struct {
	base *baseClient
}

func newAssetService(base *baseClient) AssetService {
	return &assetServicer{
		base: base,
	}
}

func (s *assetServicer) Get(options *AssetOptions) ([]*Asset, error) {
	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	ts, err := parseAssets(resp, options)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s *assetServicer) Sync(db database.SQLClient) error {
	//TODO
	return nil
}

func parseAssets(resp *http.Response, options *AssetOptions) ([]*Asset, error) {
	var target []map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	for _, route := range target {
		if route["path"] == options.Route {
			assetResponse, ok := route["assets"].(*AssetResponse)
			if !ok {
				return nil, errors.New("Error decoding asset response")
			}
			return assetResponse.Assets, nil
		}
	}
	return nil, fmt.Errorf("Error decoding asset response, Route '%s' not found", options.Route)
}
