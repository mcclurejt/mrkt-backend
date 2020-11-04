package iexcloud

import (
	"context"
	"fmt"
)

const (
	chartEndpointURL      = "/stock/%s/chart/%s/%s"
	chartBatchEndpointURL = "/stock/market/batch"
)

type ChartRange string

func (c ChartRange) String() string {
	return EnumToString(c)
}

const (
	ChartRangeMax     ChartRange = "max"
	ChartRange5y      ChartRange = "5y"
	ChartRange2y      ChartRange = "2y"
	ChartRange1y      ChartRange = "1y"
	ChartRangeYTD     ChartRange = "ytd"
	ChartRange6m      ChartRange = "6m"
	ChartRange3m      ChartRange = "3m"
	ChartRange1m      ChartRange = "1m"
	ChartRange1mm     ChartRange = "1mm"
	ChartRange5d      ChartRange = "5d"
	ChartRange5dm     ChartRange = "5dm"
	ChartRangeDate    ChartRange = "date"
	ChartRangeDynamic ChartRange = "dynamic"
)

var chartValidRanges = map[string]bool{
	ChartRangeMax.String():     true,
	ChartRange5y.String():      true,
	ChartRange2y.String():      true,
	ChartRange1y.String():      true,
	ChartRangeYTD.String():     true,
	ChartRange6m.String():      true,
	ChartRange3m.String():      true,
	ChartRange1m.String():      true,
	ChartRange1mm.String():     true,
	ChartRange5d.String():      true,
	ChartRange5dm.String():     true,
	ChartRangeDate.String():    true,
	ChartRangeDynamic.String(): true,
}

type ChartService interface {
	Get(ctx context.Context, symbol string, rang ChartRange, date string, options *ChartOptions) ([]OHLCV, error)
	GetSingleDay(ctx context.Context, symbol string, date string) ([]OHLCV, error)
	GetBatch(ctx context.Context, symbols []string, options *ChartOptions) ([]OHLCV, error)
	GetBatchSingleDay(ctx context.Context, symbols []string, date string) ([]OHLCV, error)
}

type ChartServiceOp struct {
	client *IEXCloudClient
}

var _ ChartService = &ChartServiceOp{}

type OHLCV struct {
	Symbol        string
	Date          string  `json:"date"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	Volume        int64   `json:"volume"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

type ChartOptions struct {
	ChartCloseOnly  bool       `url:"chartCloseOnly,omitempty"`
	ChartByDay      bool       `url:"chartByDay,omitempty"`
	ChartSimplify   bool       `url:"chartSimplify,omitempty"`
	ChartInterval   int        `url:"chartInterval,omitempty"`
	ChangeFromClose bool       `url:"changeFromClose,omitempty"`
	ChartLast       int        `url:"chartLast,omitempty"`
	Range           ChartRange `url:"range,omitempty"`
	ExactDate       string     `url:"exactDate,omitempty"`
	Sort            string     `url:"sort,omitempty"`
	IncludeToday    bool       `url:"includeToday,omitempty"`
}

func (c *ChartServiceOp) Get(ctx context.Context, symbol string, rang ChartRange, date string, options *ChartOptions) ([]OHLCV, error) {
	ohlcvs := []OHLCV{}
	endpoint := fmt.Sprintf(chartEndpointURL, symbol, rang, date)
	endpoint, err := c.client.addOptions(endpoint, options)
	if err != nil {
		return ohlcvs, err
	}
	err = c.client.GetJSON(ctx, endpoint, &ohlcvs)
	return ohlcvs, err
}

func (c *ChartServiceOp) GetSingleDay(ctx context.Context, symbol string, date string) ([]OHLCV, error) {
	options := &ChartOptions{
		ChartByDay: true,
	}
	return c.Get(ctx, symbol, ChartRangeDate, date, options)
}

func (c *ChartServiceOp) GetBatch(ctx context.Context, symbols []string, options *ChartOptions) ([]OHLCV, error) {
	batch := map[string]map[string][]OHLCV{}
	endpoint := chartBatchEndpointURL
	batchOptions := &BatchOptions{Symbols: SliceToString(symbols, nil), Types: "chart"}
	endpoint, err := c.client.addOptions(endpoint, batchOptions)
	if err != nil {
		return []OHLCV{}, err
	}
	endpoint, err = c.client.addOptions(endpoint, options)
	if err != nil {
		return []OHLCV{}, err
	}
	err = c.client.GetJSON(ctx, endpoint, &batch)
	ohlcvs := []OHLCV{}
	for symbol, v := range batch {
		ohlcvList := v["chart"]
		for _, ohlcv := range ohlcvList {
			ohlcv.Symbol = symbol
			ohlcvs = append(ohlcvs, ohlcv)
		}
	}
	return ohlcvs, nil
}

func (c *ChartServiceOp) GetBatchSingleDay(ctx context.Context, symbols []string, date string) ([]OHLCV, error) {
	options := &ChartOptions{
		ChartByDay: true,
		Range:      ChartRangeDate,
		ExactDate:  date,
	}
	return c.GetBatch(ctx, symbols, options)
}

// https://sandbox.iexapis.com/stable/stock/market/batch\?symbols\=aapl,fb\&types\=chart\&range\=date\&exactDate\=20201103\&chartByDay\=true\&\&token\=Tsk_3977d776bc524cada37ad2c53378f18d
