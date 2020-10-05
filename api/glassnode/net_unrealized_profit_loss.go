package glassnode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	NET_UNREALIZED_PROFIT_LOSS_FUNCTION   = "net_unrealized_profit_loss"
	NET_UNREALIZED_PROFIT_LOSS_TABLE_NAME = "NetUnrealizedProfitLoss"
)

type NetUnrealizedProfitLoss struct {
	TimeSeries []*NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossResponse struct {
	TimeSeries []*NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossEntry struct {
	Timestamp int64   `json:"t"`
	Value     float64 `json:"v"`
}

type NetUnrealizedProfitLossService interface {
	Get(coin string, interval string) (*NetUnrealizedProfitLoss, error)
	Sync(coin string, db database.SQLClient) error
}

type netUnrealizedProfitLossServiceOptions struct {
	Coin     string
	Interval string
	// since
	// until
}

func newNetUnrealizedProfitLossServiceOptions(coin string, interval string) *netUnrealizedProfitLossServiceOptions {
	return &netUnrealizedProfitLossServiceOptions{Coin: coin, Interval: interval}
}

func (o netUnrealizedProfitLossServiceOptions) ToQueryString() string {
	return fmt.Sprintf("indicators/%s?a=%s&i=%s", NET_UNREALIZED_PROFIT_LOSS_FUNCTION, o.Coin, o.Interval)
}

type netUnrealizedProfitLossServicer struct {
	base *baseClient
}

func newNetUnrealizedProfitLossService(base *baseClient) NetUnrealizedProfitLossService {
	return &netUnrealizedProfitLossServicer{
		base: base,
	}
}

func (n *netUnrealizedProfitLossServicer) Get(coin string, interval string) (*NetUnrealizedProfitLoss, error) {
	options := newNetUnrealizedProfitLossServiceOptions(coin, interval)
	resp, err := n.base.call(options)
	if err != nil {
		return nil, err
	}
	ns, err := parseNetUnrealizedProfitLoss(resp)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (n *netUnrealizedProfitLossServicer) Sync(coin string, db database.SQLClient) error {
	return nil
}

func parseNetUnrealizedProfitLoss(resp *http.Response) (*NetUnrealizedProfitLoss, error) {
	target := &[]NetUnrealizedProfitLossEntry{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	timeSeries := target
	netUnrealizedProfitLossEntries := make([]*NetUnrealizedProfitLossEntry, len(*timeSeries))
	for i, v := range *timeSeries {
		netUnrealizedProfitLossEntries[i] = &v
	}

	return &NetUnrealizedProfitLoss{TimeSeries: netUnrealizedProfitLossEntries}, nil
}
