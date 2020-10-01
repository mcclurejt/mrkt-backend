package api

import (
	"fmt"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	NET_UNREALIZED_PROFIT_LOSS_FUNCTION   = "net_unrealized_profit_loss"
	NET_UNREALIZED_PROFIT_LOSS_TABLE_NAME = "NetUnrealizedProfitLoss"
)

var (
	NET_UNREALIZED_PROFIT_LOSS_HEADERS = []string{
		"coin",
		"timestamp",
		"value",
	}
	NET_UNREALIZED_PROFIT_LOSS_COLUMNS = []string{
		"coin CHAR",
		"timestamp INT NOT NULL",
		"value FLOAT NOT NULL",
		"FOREIGN KEY (coin) REFERENCES Coin(coin)",
		"PRIMARY KEY (coin, date)",
	}
)

type NetUnrealizedProfitLoss struct {
	TimeSeries []NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossResponse struct {
	TimeSeries []NetUnrealizedProfitLossEntry
}

type NetUnrealizedProfitLossEntry struct {
	Timestamp int32 `json:t`
	Value     int32 `json:v`
}

type NetUnrealizedProfitLossService interface {
	GetTableName() string
	GetTableColumns() []string
	Get(symbol string) (NetUnrealizedProfitLoss, error)
	Insert(n NetUnrealizedProfitLoss, db database.SQLClient) error
	Sync(symbol string, db database.SQLClient) error
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
	return fmt.Sprintf("/indicators/%s?a=%s&i=%s", NET_UNREALIZED_PROFIT_LOSS_FUNCTION, o.coin, o.interval)
}
