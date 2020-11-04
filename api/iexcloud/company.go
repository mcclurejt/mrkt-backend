package iexcloud

import (
	"context"
	"fmt"
	"net/url"
)

type CompanyService interface {
	Get(ctx context.Context, symbol string) (*Company, error)
}

type CompanyServiceOp struct {
	client *IexCloudClient
}

var _ CompanyService = &CompanyServiceOp{}

type Company struct {
	Symbol         string   `json:"symbol"`
	CompanyName    string   `json:"companyName"`
	Exchange       string   `json:"exchange"`
	Industry       string   `json:"industry"`
	Website        string   `json:"website"`
	Description    string   `json:"description"`
	CEO            string   `json:"CEO"`
	SecurityName   string   `json:"securityName"`
	IssueType      string   `json:"issueType"`
	Sector         string   `json:"sector"`
	PrimarySICCode int      `json:"primarySicCode"`
	Employees      int      `json:"employees"`
	Tags           []string `json:"tags"`
	Address        string   `json:"address"`
	Address2       string   `json:"address2"`
	State          string   `json:"state"`
	City           string   `json:"city"`
	Zip            string   `json:"zip"`
	Country        string   `json:"country"`
	Phone          string   `json:"phone"`
}

func (s *CompanyServiceOp) Get(ctx context.Context, symbol string) (*Company, error) {
	company := new(Company)
	endpoint := fmt.Sprintf("/stock/%s/company", url.PathEscape(symbol))
	err := s.client.GetJSON(ctx, endpoint, company)
	return company, err
}
