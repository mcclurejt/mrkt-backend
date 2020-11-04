package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type InsiderSummaryService interface {
	Get(ctx context.Context, symbol string) ([]InsiderSummary, error)
}

type InsiderSummaryServiceOp struct {
	client *IexCloudClient
}

var _ InsiderSummaryService = &InsiderSummaryServiceOp{}

type InsiderSummary struct {
	FullName       string `json:"fullName"`
	NetTransaction int    `json:"netTransaction"`
	ReportedTitle  string `json:"reportedTitle"`
	TotalBought    int    `json:"totalBought"`
	TotalSold      int    `json:"totalSold"`
}

func (s *InsiderSummaryServiceOp) Get(ctx context.Context, symbol string) ([]InsiderSummary, error) {
	is := []InsiderSummary{}
	endpoint := fmt.Sprintf("/stock/%s/insider-summary", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &is)
	return is, err
}
