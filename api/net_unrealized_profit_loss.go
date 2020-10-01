package api

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

var (
	NET_UNREALIZED_PROFIT_LOSS_HEADERS = []string{
		"cid",
		"timestamp",
		"value",
	}
	NET_UNREALIZED_PROFIT_LOSS_COLUMNS = []string{
		"cid INT",
		"timestamp INT NOT NULL",
		"value FLOAT NOT NULL",
		"FOREIGN KEY (cid) REFERENCES Coin(cid)",
		"PRIMARY KEY (cid, timestamp)",
	}
)

type NetUnrealizedProfitLoss struct {
	TimeSeries []NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossResponse struct {
	TimeSeries []NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossEntry struct {
	Timestamp int64   `json:"t"`
	Value     float64 `json:"v"`
}

type NetUnrealizedProfitLossService interface {
	GetTableName() string
	GetTableColumns() []string
	Get(coin string, interval string) (NetUnrealizedProfitLoss, error)
	Insert(coin string, n NetUnrealizedProfitLoss, db database.SQLClient) error
	Sync(coin string, db database.SQLClient) error
}

type netUnrealizedProfitLossServiceOptions struct {
	Coin     string
	Interval string
	// since
	// until
}

func newNetUnrealizedProfitLossServiceOptions(coin string, interval string) netUnrealizedProfitLossServiceOptions {
	return netUnrealizedProfitLossServiceOptions{Coin: coin, Interval: interval}
}

func (o netUnrealizedProfitLossServiceOptions) ToQueryString() string {
	return fmt.Sprintf("indicators/%s?a=%s&i=%s", NET_UNREALIZED_PROFIT_LOSS_FUNCTION, o.Coin, o.Interval)
}

type netUnrealizedProfitLossServicer struct {
	base baseClient
}

func newNetUnrealizedProfitLossService(base baseClient) NetUnrealizedProfitLossService {
	return netUnrealizedProfitLossServicer{
		base: base,
	}
}

func (n netUnrealizedProfitLossServicer) GetTableName() string {
	return NET_UNREALIZED_PROFIT_LOSS_TABLE_NAME
}

func (n netUnrealizedProfitLossServicer) GetTableColumns() []string {
	return NET_UNREALIZED_PROFIT_LOSS_COLUMNS
}

func (n netUnrealizedProfitLossServicer) Get(coin string, interval string) (NetUnrealizedProfitLoss, error) {
	options := newNetUnrealizedProfitLossServiceOptions(coin, interval)
	resp, err := n.base.call(options)
	if err != nil {
		return NetUnrealizedProfitLoss{}, err
	}
	ns, err := parseNetUnrealizedProfitLoss(resp)
	if err != nil {
		return NetUnrealizedProfitLoss{}, err
	}
	return ns, nil
}

func (n netUnrealizedProfitLossServicer) Insert(coin string, ns NetUnrealizedProfitLoss, db database.SQLClient) error {
	coinId, err := db.GetCoinID(coin)
	if err != nil {
		return err
	}
	values := make([]interface{}, 0)
	for _, v := range ns.TimeSeries {
		values = append(values, coinId)
		values = append(values, v.Timestamp)
		values = append(values, v.Value)
	}
	return db.Insert(n.GetTableName(), NET_UNREALIZED_PROFIT_LOSS_HEADERS, values)
}

func (n netUnrealizedProfitLossServicer) Sync(coin string, db database.SQLClient) error {
	return nil
}

func parseNetUnrealizedProfitLoss(resp *http.Response) (NetUnrealizedProfitLoss, error) {
	target := &[]NetUnrealizedProfitLossEntry{}
	err := json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return NetUnrealizedProfitLoss{}, err
	}

	timeSeries := target

	netUnrealizedProfitLossEntries := make([]NetUnrealizedProfitLossEntry, len(*timeSeries))
	for i, v := range *timeSeries {
		entry := v
		netUnrealizedProfitLossEntries[i] = entry
	}

	return NetUnrealizedProfitLoss{TimeSeries: netUnrealizedProfitLossEntries}, nil
}
