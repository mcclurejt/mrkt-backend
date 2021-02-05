package iex

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ChartService interface {
	GetChart(ctx context.Context, symbol string, rang ChartRange, options *ChartOptions) ([]*OHLCV, error)
	GetChartSingleDay(ctx context.Context, symbol string, date time.Time) ([]*OHLCV, error)
}

type RealChartService struct {
	client *RealClient
}

type OHLCV struct {
	Symbol        string  `json:"symbol" attributetype:"S" keytype:"RANGE"`
	Date          string  `json:"date" attributetype:"S" keytype:"HASH"`
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

const chartEndpointURL = "/stock/%s/chart/%s"

func (c *RealClient) GetChart(ctx context.Context, symbol string, rang ChartRange, options *ChartOptions) ([]*OHLCV, error) {
	// Generate the url
	path := fmt.Sprintf(chartEndpointURL, symbol, rang)
	// Execute the request
	c.Log.WithFields(logrus.Fields{"symbol": symbol, "options": PrettyPrintStruct(*options)}).Info("Getting Charts")
	ohlcvs := []*OHLCV{}
	err := c.getWithParams(ctx, path, &ohlcvs, options)
	if err != nil {
		return nil, err
	}
	c.Log.WithFields(logrus.Fields{"symbol": symbol, "numCharts": len(ohlcvs)}).Info("Charts Received")
	// Add the symbol field to the returned objects
	for i := 0; i < len(ohlcvs); i++ {
		ohlcvs[i].Symbol = symbol
	}
	// Fill empty change percents
	ohlcvs = fillChangePercent(ohlcvs)
	return ohlcvs, nil
}

const chartEndpointURLSingleDay = "/stock/%s/chart/date/%s"

func (c *RealClient) GetChartSingleDay(ctx context.Context, symbol string, date time.Time) ([]*OHLCV, error) {
	// Generate the url
	dateString := formatDateChartRequest(date)
	path := fmt.Sprintf(chartEndpointURLSingleDay, symbol, dateString)
	// Execute the request
	c.Log.WithFields(logrus.Fields{"symbol": symbol, "date": formateDateReadable(date)}).Info("Getting Single Day Chart")
	ohlcvs := []*OHLCV{}
	options := &ChartOptions{ChartByDay: true}
	err := c.getWithParams(ctx, path, &ohlcvs, options)
	if err != nil {
		return nil, err
	}
	c.Log.WithFields(logrus.Fields{"symbol": symbol, "date": formateDateReadable(date), "numCharts": len(ohlcvs)}).Info("Charts Received")
	// Add the symbol field to the returned objects
	for i := 0; i < len(ohlcvs); i++ {
		ohlcvs[i].Symbol = symbol
	}
	// Fill empty change percents
	ohlcvs = fillChangePercent(ohlcvs)
	return ohlcvs, nil
}

// ChartDateYesterday - Returns yesterdays date as a time.Time object, useful for retrieving yesterday's chart
func ChartDateYesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}

func formatDateChartRequest(date time.Time) string {
	return date.Format("20060102")
}

func formateDateReadable(date time.Time) string {
	return date.Format("01-02-2006")
}

func fillChangePercent(ohlcvs []*OHLCV) []*OHLCV {
	for _, candle := range ohlcvs {
		if candle.ChangePercent == 0 {
			candle.ChangePercent = (candle.Close - candle.Open) / candle.Open
		}
	}
	return ohlcvs
}
