package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type InsiderTransactionsService interface {
	Get(ctx context.Context, symbol string) ([]InsiderTransaction, error)
}

type InsiderTransactionsServiceOp struct {
	client *IexCloudClient
}

var _ InsiderTransactionsService = &InsiderTransactionsServiceOp{}

type InsiderTransaction struct {
	EffectiveDate int     `json:"effectiveDate"`
	FullName      string  `json:"fullName"`
	ReportedTitle string  `json:"reportedTitle"`
	Price         float64 `json:"tranPrice"`
	Shares        int     `json:"tranShares"`
	Value         float64 `json:"tranValue"`
}

func (s *InsiderTransactionsServiceOp) Get(ctx context.Context, symbol string) ([]InsiderTransaction, error) {
	it := []InsiderTransaction{}
	endpoint := fmt.Sprintf("/stock/%s/insider-transactions", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &it)
	return it, err
}
