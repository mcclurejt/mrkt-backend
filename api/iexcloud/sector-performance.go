package iexcloud

import (
	"context"
	"fmt"
)

type SectorPerformanceService interface {
	Get(ctx context.Context) ([]SectorPerformance, error)
}

type SectorPerformanceServiceOp struct {
	client *IexCloudClient
}

var _ SectorPerformanceService = &SectorPerformanceServiceOp{}

type SectorPerformance struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Performance float64 `json:"performance"`
	LastUpdated int     `json:"lastUpdated"`
}

func (s *SectorPerformanceServiceOp) Get(ctx context.Context) ([]SectorPerformance, error) {
	sp := []SectorPerformance{}
	endpoint := fmt.Sprintf("/stock/market/sector-performance")
	err := s.client.GetJSON(ctx, endpoint, &sp)
	return sp, err
}
