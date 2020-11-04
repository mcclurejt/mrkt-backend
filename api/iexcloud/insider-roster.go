package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type InsiderRosterService interface {
	Get(ctx context.Context, symbol string) ([]InsiderRoster, error)
}

type InsiderRosterServiceOp struct {
	client *IexCloudClient
}

var _ InsiderRosterService = &InsiderRosterServiceOp{}

type InsiderRoster struct {
	EntityName string `json:"entityName"`
	Position   int    `json:"position"`
	ReportDate int    `json:"reportDate"`
}

func (s *InsiderRosterServiceOp) Get(ctx context.Context, symbol string) ([]InsiderRoster, error) {
	ir := []InsiderRoster{}
	endpoint := fmt.Sprintf("/stock/%s/insider-roster", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &ir)
	return ir, err
}
