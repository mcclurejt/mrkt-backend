package iexcloud

import (
	"context"
)

type StatusService interface {
	Get(ctx context.Context) (*Status, error)
}

type StatusServiceOp struct {
	client *IEXCloudClient
}

var _ StatusService = &StatusServiceOp{}

// Status models the IEX Cloud API system status
type Status struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    int    `json:"time"`
}

func (s *StatusServiceOp) Get(ctx context.Context) (*Status, error) {
	status := new(Status)
	endpoint := "/status"
	err := s.client.GetJSON(ctx, endpoint, status)
	return status, err
}
