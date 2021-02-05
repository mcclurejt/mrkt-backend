package iex

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// common - Models shared across the iexcloud API

const DefaultTimeStampFormat = time.RFC3339

type BidAsk struct {
	Price     float64 `json:"price"`
	Size      int     `json:"size"`
	Timestamp int     `json:"timestamp"`
}

// Quote models the data returned from the IEX Cloud /quote endpoint.
type Quote struct {
	Symbol                 string  `json:"symbol"`
	CompanyName            string  `json:"companyName"`
	Exchange               string  `json:"primaryExchange"`
	CalculationPrice       string  `json:"calculationPrice"`
	Open                   float64 `json:"open"`
	OpenTime               int     `json:"openTime"`
	OpenSource             string  `json:"openSource"`
	Close                  float64 `json:"close"`
	CloseTime              int     `json:"closeTime"`
	CloseSource            string  `json:"closeSource"`
	High                   float64 `json:"high"`
	Low                    float64 `json:"low"`
	LowTime                float64 `json:"lowTime"`
	LowSource              string  `json:"lowSource"`
	LatestPrice            float64 `json:"latestPrice"`
	LatestSource           string  `json:"latestSource"`
	LatestTime             string  `json:"latestTime"`
	LatestUpdate           int     `json:"latestUpdate"`
	LatestVolume           int     `json:"latestVolume"`
	IexRealtimePrice       float64 `json:"iexRealtimePrice"`
	IexRealtimeSize        int     `json:"iexRealtimeSize"`
	IexLastUpdated         int     `json:"iexLastUpdated"`
	DelayedPrice           float64 `json:"delayedPrice"`
	DelayedPriceTime       int     `json:"delayedPriceTime"`
	OddLotDelayedPrice     int     `json:"oddLotDelayedPrice"`
	OddLotDelayedPriceTime int     `json:"oddLotDelayedPriceTime"`
	ExtendedPrice          float64 `json:"extendedPrice"`
	ExtendedChange         float64 `json:"extendedChange"`
	ExtendedChangePercent  float64 `json:"extendedChangePercent"`
	ExtendedPriceTime      int     `json:"extendedPriceTime"`
	PreviousClose          float64 `json:"previousClose"`
	Change                 float64 `json:"change"`
	ChangePercent          float64 `json:"changePercent"`
	Volume                 int     `json:"volume"`
	IexMarketPercent       float64 `json:"iexMarketPercent"`
	IexVolume              int     `json:"iexVolume"`
	AvgTotalVolume         int     `json:"avgTotalVolume"`
	IexBidPrice            float64 `json:"iexBidPrice"`
	IexBidSize             int     `json:"iexBidSize"`
	IexAskPrice            float64 `json:"iexAskPrice"`
	IexAskSize             int     `json:"iexAskSize"`
	IexOpen                float64 `json:"iexOpen"`
	IexOpenTime            int     `json:"iexOpenTime"`
	IexClose               float64 `json:"iexClose"`
	IexCloseTime           int     `json:"iexCloseTime"`
	MarketCap              int     `json:"marketCap"`
	PERatio                float64 `json:"peRatio"`
	Week52High             float64 `json:"week52High"`
	Week52Low              float64 `json:"week52Low"`
	YTDChange              float64 `json:"ytdChange"`
	LastTradeTime          int     `json:"lastTradeTime"`
}

// Trade models a trade for a quote.
type Trade struct {
	Price                 float64 `json:"price"`
	Size                  int     `json:"size"`
	TradeID               int     `json:"tradeId"`
	IsISO                 bool    `json:"isISO"`
	IsOddLot              bool    `json:"isOddLot"`
	IsOutsideRegularHours bool    `json:"isOutsideRegularHours"`
	IsSinglePriceCross    bool    `json:"isSinglePriceCross"`
	IsTradeThroughExempt  bool    `json:"isTradeThroughExempt"`
	Timestamp             int     `json:"timestamp"`
}

// StrToPtr - returns a pointer to the provided string since &"" is not allowed by go
func StrToPtr(s string) *string {
	return &s
}

// SliceToString - takes a slice of string-like objects and converts them to a string containing the items separated by the separator (comma is default)
func SliceToString(arr interface{}, sep *string) string {
	t := reflect.TypeOf(arr)
	if t.Kind() != reflect.Slice {
		panic(arr)
	}
	if sep == nil {
		s := ","
		sep = &s
	}
	v := reflect.ValueOf(arr)
	l := v.Len()
	stringArr := make([]string, l)
	for i := 0; i < l; i++ {
		entry := v.Index(i)
		stringArr[i] = entry.String()
	}
	return url.PathEscape(strings.Join(stringArr, *sep))
}

// EnumToString - Takes a custom-typed object and converts it to a string
func EnumToString(e interface{}) string {
	return reflect.ValueOf(e).String()
}

// DateToTimestamp - Takes a date in the format of yyyy-mm-dd and converts it to a timestamp
func DateToTimestamp(d string) (string, error) {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		return "", err
	}
	return t.Format(DefaultTimeStampFormat), nil
}

// TimeToTimestamp - Takes a Unix time eg: "1257894000" and converts it to a timestamp
func TimeToTimestamp(seconds int64) string {
	t := time.Unix(seconds, 0)
	return t.Format(DefaultTimeStampFormat)
}

func PrettyPrintStruct(v interface{}) string {
	output := ""
	rv := reflect.ValueOf(v)
	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsZero() {
			output += fmt.Sprintf("%s: %v, ", rv.Type().Field(i).Name, rv.Field(i).Interface())
		}
	}
	return output
}
