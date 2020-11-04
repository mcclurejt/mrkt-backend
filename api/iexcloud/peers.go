package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type PeersService interface {
	Get(ctx context.Context, symbol string) ([]string, error)
}

type PeersServiceOp struct {
	client *IEXCloudClient
}

var _ PeersService = &PeersServiceOp{}

func (s *PeersServiceOp) Get(ctx context.Context, symbol string) ([]string, error) {
	peers := []string{}
	endpoint := fmt.Sprintf("/stock/%s/peers", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &peers)
	return peers, err
}
