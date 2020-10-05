package alphavantage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/mcclurejt/mrkt-backend/api/dynamodb"
)

const (
	COMPANY_OVERVIEW_FUNCTION   = "OVERVIEW"
	COMPANY_OVERVIEW_TABLE_NAME = "CompanyOverview"
)

type CompanyOverview struct {
	Name                       string  `json:"Name"`
	AssetType                  string  `json:"AssetType"`
	Description                string  `json:"Description"`
	Symbol                     string  `json:"Symbol"`
	Exchange                   string  `json:"Exchange"`
	Currency                   string  `json:"Currency"`
	Country                    string  `json:"Country"`
	Sector                     string  `json:"Sector"`
	Industry                   string  `json:"Industry"`
	Address                    string  `json:"Address"`
	FullTimeEmployees          string  `json:"FullTimeEmployees"`
	FiscalYearEnd              string  `json:"FiscalYearEnd"`
	LatestQuarter              string  `json:"LatestQuarter"`
	MarketCapitalization       float64 `json:"MarketCapitalization"`
	EBITDA                     int64   `json:"EBITDA"`
	PERatio                    float64 `json:"PERatio"`
	PEGRatio                   float64 `json:"PEGRatio"`
	BookValue                  float64 `json:"BookValue"`
	DividendPerShare           float64 `json:"DividendPerShare"`
	DividendYield              float64 `json:"DividendYield"`
	EPS                        float64 `json:"EPS"`
	RevenuePerShareTTM         float64 `json:"RevenuePerShareTTM"`
	ProfitMargin               float64 `json:"ProfitMargin"`
	OperatingMarginTTM         float64 `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          float64 `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          float64 `json:"ReturnOnEquityTTM"`
	RevenueTTM                 int64   `json:"RevenueTTM"`
	GrossProfitTTM             int64   `json:"GrossProfitTTM"`
	DilutedEPSTTM              float64 `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY float64 `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  float64 `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         float64 `json:"AnalystTargetPrice"`
	TrailingPE                 float64 `json:"TrailingPE"`
	ForwardPE                  float64 `json:"ForwardPE"`
	PriceToSalesRatioTTM       float64 `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           float64 `json:"PriceToBookRatio"`
	EVToRevenue                float64 `json:"EVToRevenue"`
	EvToEBITDA                 float64 `json:"EvToEBITDA"`
	Beta                       float64 `json:"Beta"`
	FiftyTwoWeekHigh           float64 `json:"52WeekHigh"`
	FiftyTwoWeekLow            float64 `json:"52WeekLow"`
	FiftyDayMovingAverage      float64 `json:"50DayMovingAverage"`
	TwoHundredDayMovingAverage float64 `json:"200DayMovingAverage"`
	SharesOutstanding          int64   `json:"SharesOutstanding"`
	SharesFloat                int64   `json:"SharesFloat"`
	SharesShort                int64   `json:"SharesShort"`
	SharesShortPriorMonth      int64   `json:"SharesShortPriorMonth"`
	ShortRatio                 float64 `json:"ShortRatio"`
	ShortPercentOutstanding    float64 `json:"ShortPercentOutstanding"`
	ShortPercentFloat          float64 `json:"ShortPercentFloat"`
	PercentInsiders            float64 `json:"PercentInsiders"`
	PercentInstitutions        float64 `json:"PercentInstitutions"`
	ForwardAnnualDividendRate  float64 `json:"ForwardAnnualDividendRate"`
	ForwardAnnualDividendYield float64 `json:"ForwardAnnualDividendYield"`
	PayoutRatio                float64 `json:"PayoutRatio"`
	DividendDate               string  `json:"DividendDate"`
	ExDividendDate             string  `json:"ExDividendDate"`
	LastSplitFactor            string  `json:"LastSplitFactor"`
	LastSplitDate              string  `json:"LastSplitDate"`
}

func (co *CompanyOverview) UnmarshalJSON(b []byte) error {
	var src map[string]interface{}

	err := json.Unmarshal(b, &src)
	if err != nil {
		return err
	}
	target := &CompanyOverview{}
	dest := reflect.ValueOf(target).Elem()
	destType := reflect.TypeOf(target).Elem()
	for i := 0; i < dest.NumField(); i++ {
		field := dest.Field(i)
		tag := destType.Field(i).Tag.Get("json")
		val, ok := src[tag]
		valString, stringOK := val.(string)
		if ok && stringOK && valString != "None" {
			switch field.Interface().(type) {
			case int, int32, int64:
				out, _ := strconv.ParseInt(val.(string), 10, 64)
				field.SetInt(out)
				break
			case float32, float64:
				out, _ := strconv.ParseFloat(val.(string), 64)
				field.SetFloat(out)
				break
			case string:
				field.SetString(val.(string))
				break
			}
		}
	}
	*co = CompanyOverview(*target)
	return nil
}

type CompanyOverviewService interface {
	GetCreateTableInput() *dynamodb.CreateTableInput
	GetPutItemInput() *dynamodb.PutItemInput

	Get(*CompanyOverviewOptions) (*CompanyOverview, error)
	Sync(symbol string, db db.Client) error
}

type CompanyOverviewOptions struct {
	Symbol string
}

func newCompanyOverviewOptions(symbol string) *CompanyOverviewOptions {
	return &CompanyOverviewOptions{Symbol: symbol}
}

func (o *CompanyOverviewOptions) ToQueryString() string {
	return fmt.Sprintf("&function=%s&symbol=%s", COMPANY_OVERVIEW_FUNCTION, o.Symbol)
}

type companyOverviewServicer struct {
	base *baseClient
}

func newCompanyOverviewService(base *baseClient) CompanyOverviewService {
	return &companyOverviewServicer{
		base: base,
	}
}

func (s *companyOverviewServicer) GetCreateTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Symbol"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Symbol"),
				KeyType:       aws.String("HASH"),
			},
		},
		BillingMode: aws.String(db.DefaultBillingMode),
		TableName:   aws.String(COMPANY_OVERVIEW_TABLE_NAME),
	}
}

func (s *companyOverviewServicer) GetPutItemInput() *dynamodb.PutItemInput {
	return &dynamodb.PutItemInput{
		TableName: aws.String(COMPANY_OVERVIEW_TABLE_NAME),
	}
}

func (s *companyOverviewServicer) Get(options *CompanyOverviewOptions) (*CompanyOverview, error) {
	if options.Symbol == "" {
		return nil, errors.New("Must provide a symbol")
	}

	resp, err := s.base.call(options)
	if err != nil {
		return nil, err
	}

	co, err := parseCompanyOverview(resp)
	if err != nil {
		_, ok := err.(*AlphaVantageRateExceededError)
		if ok {
			time.Sleep(defaultRetryPeriod * time.Second)
			return s.Get(options)
		}
		return nil, err
	}

	return co, nil
}

func (s *companyOverviewServicer) Sync(symbol string, db db.Client) error {
	// TODO
	return nil
}

func parseCompanyOverview(resp *http.Response) (*CompanyOverview, error) {
	target := &CompanyOverview{}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	err := json.Unmarshal(body, &target)
	if err != nil {
		return nil, err
	}

	// turn FiscalYearEnd to a date
	monthString := target.FiscalYearEnd
	monthNumber := 1
	for monthString != time.Month(monthNumber).String() && monthNumber <= 12 {
		monthNumber += 1
	}
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	t := time.Date(time.Now().Year(), time.Month(monthNumber), 1, 0, 0, 0, 0, loc)
	if err != nil {
		return nil, err
	}
	target.FiscalYearEnd = fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())

	// check to see if the rate was exceeded and no objects were returned (still gives 200 status code)
	if target.Name == "" {
		return nil, &AlphaVantageRateExceededError{}
	}

	return target, nil
}
