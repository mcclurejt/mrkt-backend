package glassnode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	BTC = "BTC"
	ETH = "ETH"
)

const (
	COIN_TABLE_NAME = "Coin"
)

type CoinResponse struct {
	Data []*CoinEntry `json:data`
}

type CoinEntry struct {
	Symbol string `json:symbol`
}

type CoinService interface {
	Get() ([]*CoinEntry, error)
	Sync(db database.SQLClient) error
}

type coinServiceOptions struct {
}

func newCoinServiceOptions() *coinServiceOptions {
	return &coinServiceOptions{}
}

func (o *coinServiceOptions) ToQueryString() string {
	return fmt.Sprintf("")
}

type coinServicer struct {
	base *baseClient
}

func newCoinService(base *baseClient) CoinService {
	return &coinServicer{
		base: base,
	}
}

func (s *coinServicer) Get() ([]*CoinEntry, error) {
	options := newCoinServiceOptions()
	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	ts, err := parseCoin(resp)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s coinServicer) Sync(db database.SQLClient) error {
	//TODO
	return nil
}

func parseCoin(resp *http.Response) ([]*CoinEntry, error) {
	target := &CoinResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	coins := make([]*CoinEntry, len(target.Data))
	for i, v := range target.Data {
		coins[i] = &CoinEntry{Symbol: v.Symbol}
	}
	return coins, nil
}
