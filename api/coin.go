package api

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

var (
	COIN_HEADERS = []string{
		"cid",
		"symbol",
	}
	COIN_INSERTION_HEADERS = []string{
		"symbol",
	}
	COIN_COLUMNS = []string{
		"cid INT NOT NULL AUTO_INCREMENT PRIMARY KEY",
		"symbol VARCHAR (8) NOT NULL UNIQUE",
	}
)

type Coins struct {
	Data []string
}

type CoinResponse struct {
	Data []CoinEntry `json:data`
}

type CoinEntry struct {
	Symbol string `json:symbol`
}

type CoinService interface {
	GetTableName() string
	GetTableColumns() []string
	Get() (Coins, error)
	Insert(c Coins, db database.SQLClient) error
	Sync(db database.SQLClient) error
}

type coinServiceOptions struct {
}

func newCoinServiceOptions() coinServiceOptions {
	return coinServiceOptions{}
}

func (o coinServiceOptions) ToQueryString() string {
	return fmt.Sprintf("")
}

type coinServicer struct {
	base baseClient
}

func newCoinService(base baseClient) CoinService {
	return coinServicer{
		base: base,
	}
}

func (s coinServicer) GetTableName() string {
	return COIN_TABLE_NAME
}

func (s coinServicer) GetTableColumns() []string {
	return COIN_COLUMNS
}

func (s coinServicer) Get() (Coins, error) {
	options := newCoinServiceOptions()
	resp, err := s.base.call(options)
	if err != nil {
		return Coins{}, err
	}

	ts, err := parseCoin(resp)
	if err != nil {
		return Coins{}, err
	}
	return ts, nil
}

func (s coinServicer) Insert(c Coins, db database.SQLClient) error {
	cs := c.Data
	values := make([]interface{}, len(cs))
	for i, v := range c.Data {
		values[i] = v
	}
	return db.Insert(s.GetTableName(), COIN_INSERTION_HEADERS, values)
}

func (s coinServicer) Sync(db database.SQLClient) error {
	//TODO
	return nil
}

func parseCoin(resp *http.Response) (Coins, error) {
	target := &CoinResponse{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return Coins{}, err
	}

	coin := make([]string, len(target.Data))
	for i, v := range target.Data {
		coin[i] = v.Symbol
	}
	return Coins{Data: coin}, nil
}
