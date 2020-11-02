package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type BookService interface {
	Get(context.Context, string) (*Book, error)
}

type BookServiceOp struct {
	client *IexCloudClient
}

var _ BookService = &BookServiceOp{}

type Book struct {
	Quote  Quote    `json:"quote"`
	Bids   []BidAsk `json:"bids"`
	Asks   []BidAsk `json:"asks"`
	Trades []Trade  `json:"trades"`
}

func (s *BookServiceOp) Get(ctx context.Context, symbol string) (*Book, error) {
	book := new(Book)
	endpoint := fmt.Sprintf("/stock/%s/book", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, &book)
	return book, err
}
