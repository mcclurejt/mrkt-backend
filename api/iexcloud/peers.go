package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type PeersService interface {
	Get(context.Context, string) (*[]string, error)
}

type PeersServiceOp struct {
	client *IexCloudClient
}

var _ PeersService = &PeersServiceOp{}

func (s *PeersServiceOp) Get(ctx context.Context, symbol string) (*[]string, error) {
	peers := new([]string)
	endpoint := fmt.Sprintf("/stock/%s/peers", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &peers)
	return peers, err
}
