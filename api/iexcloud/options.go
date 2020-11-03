package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type OptionsService interface {
	GetOptionsExpDates(context.Context, string) (*OptionsExpDates, error)
	GetOptionContracts(context.Context, string, string, string) (*[]OptionContract, error)
}

type OptionsServiceOp struct {
	client *IexCloudClient
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
	options := new(OptionsExpDates)
	endpoint := fmt.Sprintf("/stock/%s/options", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &options)
	return options, err
}

func (s *OptionsServiceOp) GetOptionContracts(ctx context.Context, symbol string, expiration string, side string) (*[]OptionContract, error) {
	options := new([]OptionContract)
	endpoint := fmt.Sprintf("/stock/%s/options/%s/%s", url.PathEscape(symbol), expiration, side)
	err := s.client.GetJSON(ctx, endpoint, &options)
	return options, err
}
