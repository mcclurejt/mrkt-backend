package glassnode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Interval string

const (
	nuplRoute                = "/v1/metrics/indicators/net_unrealized_profit_loss"
	nuplTableName            = "NetUnrealizedProfitLoss"
	Interval24h     Interval = "24h"
	Interval1h      Interval = "1h"
	IntervalDefault Interval = Interval24h
)

type NetUnrealizedProfitLossEntry struct {
	Timestamp int64   `json:"t"`
	Value     float64 `json:"v"`
}

type NetUnrealizedProfitLossService interface {
	Get(options *NetUnrealizedProfitLossOptions) ([]*NetUnrealizedProfitLossEntry, error)
	GetBatch(options *NetUnrealizedProfitLossOptions, ch chan<- ResultError)
}

type NetUnrealizedProfitLossOptions struct {
	Asset    string
	Interval Interval
	Since    *int
	Until    *int
	sync.RWMutex
}

func DefaultNetUnrealizedProfitLossOptions() *NetUnrealizedProfitLossOptions {
	return &NetUnrealizedProfitLossOptions{
		Asset:    BTC,
		Interval: IntervalDefault,
		Since:    nil,
		Until:    nil,
	}
}

func (o *NetUnrealizedProfitLossOptions) ToQueryString() string {
	o.Lock()
	defer o.Unlock()
	qs := fmt.Sprintf("%s?a=%s&i=%s", nuplRoute, o.Asset, o.Interval)
	if o.Since != nil {
		qs += fmt.Sprintf("&s=%s", *o.Since)
	}
	if o.Until != nil {
		qs += fmt.Sprintf("&u=%s", *o.Until)
	}
	return qs
}

type netUnrealizedProfitLossServicer struct {
	base *baseClient
}

func newNetUnrealizedProfitLossService(base *baseClient) NetUnrealizedProfitLossService {
	return &netUnrealizedProfitLossServicer{
		base: base,
	}
}

func (n netUnrealizedProfitLossServicer) Get(options *NetUnrealizedProfitLossOptions) ([]*NetUnrealizedProfitLossEntry, error) {
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

func (n netUnrealizedProfitLossServicer) GetBatch(options *NetUnrealizedProfitLossOptions, ch chan<- ResultError) {
	nupl, err := n.Get(options)
	if err != nil {
		ch <- ResultError{Error: err}
	} else {
		ch <- ResultError{Result: nupl}
	}
}

func parseNetUnrealizedProfitLoss(resp *http.Response) ([]*NetUnrealizedProfitLossEntry, error) {
	target := []*NetUnrealizedProfitLossEntry{}
	err := json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}
