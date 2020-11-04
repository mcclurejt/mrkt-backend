package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type OptionsSide string

const (
	OptionsSideCall OptionsSide = "call"
	OptionsSidePut  OptionsSide = "put"
)

func (s OptionsSide) String() string {
	return EnumToString(s)
}

type OptionsService interface {
	GetOptionsExpDates(ctx context.Context, symbol string) (*OptionsExpDates, error)
	GetOptionContracts(ctx context.Context, symbol string, expiration string, side OptionsSide) ([]OptionContract, error)
}

type OptionsServiceOp struct {
	client *IEXCloudClient
}

var _ OptionsService = &OptionsServiceOp{}

type OptionsExpDates []string

type OptionContract struct {
	Symbol         string  `json:"symbol"`
	Id             string  `json:"id"`
	ExpirationDate string  `json:"expirationDate"`
	ContractSize   int     `json:"contractSize"`
	StrikePrice    int     `json:"strikePrice"`
	ClosingPrice   float64 `json:"closingPrice"`
	Side           string  `json:"side"`
	Type           string  `json:"type"`
	Volume         int     `json:"volume"`
	OpenInterest   int     `json:"openInterest"`
	Bid            float64 `json:"bid"`
	Ask            float64 `json:"ask"`
	LastUpdated    string  `json:"lastUpdated"`
	IsAdjusted     bool    `json:"isAdjusted"`
}

func (s *OptionsServiceOp) GetOptionsExpDates(ctx context.Context, symbol string) (*OptionsExpDates, error) {
	options := &OptionsExpDates{}
	endpoint := fmt.Sprintf("/stock/%s/options", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, options)
	return options, err
}

func (s *OptionsServiceOp) GetOptionContracts(ctx context.Context, symbol string, expiration string, side OptionsSide) ([]OptionContract, error) {
	options := []OptionContract{}
	endpoint := fmt.Sprintf("/stock/%s/options/%s/%s", url.PathEscape(symbol), expiration, side.String())
	err := s.client.GetJSON(ctx, endpoint, &options)
	return options, err
}
