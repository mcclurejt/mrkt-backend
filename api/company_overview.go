package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/mcclurejt/mrkt-backend/database"
)

const (
	COMPANY_OVERVIEW_FUNCTION   = "OVERVIEW"
	COMPANY_OVERVIEW_TABLE_NAME = "CompanyOverview"
)

var (
	COMPANY_OVERVIEW_HEADERS = []string{
		"ID",
		"Name",
		"AssetType",
		"Description",
		"Symbol",
		"Exchange",
		"Currency",
		"Country",
		"Sector",
		"Industry",
		"Address",
		"FullTimeEmployees",
		"FiscalYearEnd",
		"LatestQuarter",
		"MarketCapitalization",
		"EBITDA",
		"PERatio",
		"PEGRatio",
		"BookValue",
		"DividendPerShare",
		"DividendYield",
		"EPS",
		"RevenuePerShareTTM",
		"ProfitMargin",
		"OperatingMarginTTM",
		"ReturnOnAssetsTTM",
		"ReturnOnEquityTTM",
		"RevenueTTM",
		"GrossProfitTTM",
		"DilutedEPSTTM",
		"QuarterlyEarningsGrowthYOY",
		"QuarterlyRevenueGrowthYOY",
		"AnalystTargetPrice",
		"TrailingPE",
		"ForwardPE",
		"PriceToSalesRatioTTM",
		"PriceToBookRatio",
		"EVToRevenue",
		"EvToEBITDA",
		"Beta",
		"FiftyTwoWeekHigh",
		"FiftyTwoWeekLow",
		"FiftyDayMovingAverage",
		"TwoHundredDayMovingAverage",
		"SharesOutstanding",
		"SharesFloat",
		"SharesShort",
		"SharesShortPriorMonth",
		"ShortRatio",
		"ShortPercentOutstanding",
		"ShortPercentFloat",
		"PercentInsiders",
		"PercentInstitutions",
		"ForwardAnnualDividendRate",
		"ForwardAnnualDividendYield",
		"PayoutRatio",
		"DividendDate",
		"ExDividendDate",
		"LastSplitFactor",
		"LastSplitDate",
	}
	COMPANY_OVERVIEW_COLUMNS = []string{
		"ID INT NOT NULL UNIQUE",
		"Name VARCHAR(64) NOT NULL",
		"AssetType VARCHAR(32)",
		"Description TEXT",
		"Symbol VARCHAR(8)",
		"Exchange VARCHAR(32)",
		"Currency VARCHAR(8)",
		"Country VARCHAR(32)",
		"Sector VARCHAR(32)",
		"Industry VARCHAR(32)",
		"Address TEXT",
		"FullTimeEmployees INT",
		"FiscalYearEnd DATE",
		"LatestQuarter DATE",
		"MarketCapitalization BIGINT",
		"EBITDA BIGINT",
		"PERatio FLOAT",
		"PEGRatio FLOAT",
		"BookValue FLOAT",
		"DividendPerShare FLOAT",
		"DividendYield FLOAT",
		"EPS FLOAT",
		"RevenuePerShareTTM FLOAT",
		"ProfitMargin FLOAT",
		"OperatingMarginTTM FLOAT",
		"ReturnOnAssetsTTM FLOAT",
		"ReturnOnEquityTTM FLOAT",
		"RevenueTTM BIGINT",
		"GrossProfitTTM BIGINT",
		"DilutedEPSTTM FLOAT",
		"QuarterlyEarningsGrowthYOY FLOAT",
		"QuarterlyRevenueGrowthYOY FLOAT",
		"AnalystTargetPrice FLOAT",
		"TrailingPE FLOAT",
		"ForwardPE FLOAT",
		"PriceToSalesRatioTTM FLOAT",
		"PriceToBookRatio FLOAT",
		"EVToRevenue FLOAT",
		"EvToEBITDA FLOAT",
		"Beta FLOAT",
		"FiftyTwoWeekHigh FLOAT",
		"FiftyTwoWeekLow FLOAT",
		"FiftyDayMovingAverage FLOAT",
		"TwoHundredDayMovingAverage FLOAT",
		"SharesOutstanding BIGINT",
		"SharesFloat BIGINT",
		"SharesShort BIGINT",
		"SharesShortPriorMonth BIGINT",
		"ShortRatio FLOAT",
		"ShortPercentOutstanding FLOAT",
		"ShortPercentFloat FLOAT",
		"PercentInsiders FLOAT",
		"PercentInstitutions FLOAT",
		"ForwardAnnualDividendRate FLOAT",
		"ForwardAnnualDividendYield FLOAT",
		"PayoutRatio FLOAT",
		"DividendDate DATE",
		"ExDividendDate DATE",
		"LastSplitFactor VARCHAR(32)",
		"LastSplitDate DATE",
		"FOREIGN KEY (ID) REFERENCES Ticker(id)",
		"PRIMARY KEY (ID)",
	}
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
	GetTableName() string
	GetTableColumns() []string
	Get(symbol string) (CompanyOverview, error)
	Insert(co CompanyOverview, db database.SQLClient) error
	Sync(symbol string, db database.SQLClient) error
}

type companyOverviewServiceOptions struct {
	Symbol string
}

func newCompanyOverviewServiceOptions(symbol string) companyOverviewServiceOptions {
	return companyOverviewServiceOptions{Symbol: symbol}
}

func (o companyOverviewServiceOptions) ToQueryString() string {
	return fmt.Sprintf("&function=%s&symbol=%s", COMPANY_OVERVIEW_FUNCTION, o.Symbol)
}

type companyOverviewServicer struct {
	base baseClient
}

func newCompanyOverviewService(base baseClient) CompanyOverviewService {
	return companyOverviewServicer{
		base: base,
	}
}

func (s companyOverviewServicer) GetTableName() string {
	return COMPANY_OVERVIEW_TABLE_NAME
}

func (s companyOverviewServicer) GetTableColumns() []string {
	return COMPANY_OVERVIEW_COLUMNS
}

func (s companyOverviewServicer) Get(symbol string) (CompanyOverview, error) {
	options := newCompanyOverviewServiceOptions(symbol)
	resp, err := s.base.call(options)
	if err != nil {
		return CompanyOverview{}, err
	}

	co, err := parseCompanyOverview(resp)
	if err != nil {
		_, ok := err.(*AlphaVantageRateExceededError)
		if ok {
			time.Sleep(DEFAULT_RETRY_PERIOD_SECONDS * time.Second)
			return s.Get(symbol)
		}
		return CompanyOverview{}, err
	}

	return co, nil
}

func (s companyOverviewServicer) Insert(co CompanyOverview, db database.SQLClient) error {
	tickerID, err := db.GetTickerID(co.Symbol)
	if err != nil {
		return err
	}

	values := make([]interface{}, len(COMPANY_OVERVIEW_HEADERS))
	ref := reflect.ValueOf(co)
	values[0] = tickerID
	for i := 1; i < len(COMPANY_OVERVIEW_HEADERS); i++ {
		s := COMPANY_OVERVIEW_HEADERS[i]
		values[i] = ref.FieldByName(s).Interface()
		if values[i] == "" {
			values[i] = nil
		}
	}
	return db.Insert(s.GetTableName(), COMPANY_OVERVIEW_HEADERS, values)
}

func (s companyOverviewServicer) Sync(symbol string, db database.SQLClient) error {
	// TODO
	return nil
}

func parseCompanyOverview(resp *http.Response) (CompanyOverview, error) {
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
		return CompanyOverview{}, err
	}

	// turn FiscalYearEnd to a date
	monthString := target.FiscalYearEnd
	monthNumber := 1
	for monthString != time.Month(monthNumber).String() && monthNumber <= 12 {
		monthNumber += 1
	}
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return CompanyOverview{}, err
	}
	t := time.Date(time.Now().Year(), time.Month(monthNumber), 1, 0, 0, 0, 0, loc)
	if err != nil {
		return CompanyOverview{}, err
	}
	target.FiscalYearEnd = fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())

	// check to see if the rate was exceeded and no objects were returned (still gives 200 status code)
	if target.Name == "" {
		return CompanyOverview{}, &AlphaVantageRateExceededError{}
	}

	return *target, nil
}