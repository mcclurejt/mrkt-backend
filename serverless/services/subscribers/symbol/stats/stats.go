package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	iex "github.com/goinvest/iexcloud/v2"

	"github.com/sirupsen/logrus"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

var (
	iexClient *iex.Client
	ddbClient dynamodbiface.DynamoDBAPI
	log       *logrus.Logger
)

func init() {
	iexClient = iex.NewClient("pk_1d8a2228abd84b0598a6cf91a5d09f63")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return
	}
	ddbClient = ddb.New(awsSession)
	log = logrus.New()
}

func processItem(item map[string]events.DynamoDBAttributeValue) error {
	// Retrieve advanced stats
	symbol, ok := item["Symbol"]
	if !ok {
		return errors.New("Symbol Key Not Found")
	}
	log.Info("Retrieving stats for %s", symbol.String())
	t := time.Now()
	data, err := iexClient.AdvancedStats(context.Background(), symbol.String())
	if err != nil {
		return err
	}
	log.Infof("Retrieved stats in %.2fs", time.Now().Sub(t).Seconds())
	// Form the request
	input := &ddb.PutItemInput{
		TableName: aws.String("Stats"),
		Item: map[string]*ddb.AttributeValue{
			"Symbol": {
				S: aws.String(symbol.String()),
			},
			"MarketCap": {
				N: aws.String(fmt.Sprintf("%f", data.MarketCap)),
			},
			"Week52High": {
				N: aws.String(fmt.Sprintf("%f", data.Week52High)),
			},
			"Week52Low": {
				N: aws.String(fmt.Sprintf("%f", data.Week52Low)),
			},
			"Week52Change": {
				N: aws.String(fmt.Sprintf("%f", data.Week52Change)),
			},
			"SharesOutstanding": {
				N: aws.String(fmt.Sprintf("%f", data.SharesOutstanding)),
			},
			"Float": {
				N: aws.String(fmt.Sprintf("%f", data.Float)),
			},
			"Avg10Volume": {
				N: aws.String(fmt.Sprintf("%f", data.Avg10Volume)),
			},
			"Avg30Volume": {
				N: aws.String(fmt.Sprintf("%f", data.Avg30Volume)),
			},
			"Day200MovingAvg": {
				N: aws.String(fmt.Sprintf("%f", data.Day200MovingAvg)),
			},
			"Day50MovingAvg": {
				N: aws.String(fmt.Sprintf("%f", data.Day50MovingAvg)),
			},
			"Employees": {
				N: aws.String(fmt.Sprintf("%d", data.Employees)),
			},
			"TTMEPS": {
				N: aws.String(fmt.Sprintf("%f", data.TTMEPS)),
			},
			"TTMDividendRate": {
				N: aws.String(fmt.Sprintf("%f", data.TTMDividendRate)),
			},
			"DividendYield": {
				N: aws.String(fmt.Sprintf("%f", data.DividendYield)),
			},
			"NextDividendDate": {
				S: aws.String(time.Time(data.NextDividendDate).Format("2006-01-02")),
			},
			"ExDividendDate": {
				S: aws.String(time.Time(data.ExDividendDate).Format("2006-01-02")),
			},
			"NextEarningsDate": {
				S: aws.String(time.Time(data.NextEarningsDate).Format("2006-01-02")),
			},
			"PERatio": {
				N: aws.String(fmt.Sprintf("%f", data.PERatio)),
			},
			"Beta": {
				N: aws.String(fmt.Sprintf("%f", data.Beta)),
			},
			"MaxChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.MaxChangePercent)),
			},
			"Year5ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Year5ChangePercent)),
			},
			"Year2ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Year2ChangePercent)),
			},
			"Year1ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Year1ChangePercent)),
			},
			"YTDChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.YTDChangePercent)),
			},
			"Month6ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Month6ChangePercent)),
			},
			"Month3ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Month3ChangePercent)),
			},
			"Month1ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Month1ChangePercent)),
			},
			"Day30ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Day30ChangePercent)),
			},
			"Day5ChangePercent": {
				N: aws.String(fmt.Sprintf("%f", data.Day5ChangePercent)),
			},
			"TotalCash": {
				N: aws.String(fmt.Sprintf("%f", data.TotalCash)),
			},
			"CurrentDebt": {
				N: aws.String(fmt.Sprintf("%f", data.CurrentDebt)),
			},
			"Revenue": {
				N: aws.String(fmt.Sprintf("%f", data.Revenue)),
			},
			"GrossProfit": {
				N: aws.String(fmt.Sprintf("%f", data.GrossProfit)),
			},
			"TotalRevenue": {
				N: aws.String(fmt.Sprintf("%f", data.TotalRevenue)),
			},
			"EBITDA": {
				N: aws.String(fmt.Sprintf("%f", data.EBITDA)),
			},
			"RevenuePerShare": {
				N: aws.String(fmt.Sprintf("%f", data.RevenuePerShare)),
			},
			"RevenuePerEmployee": {
				N: aws.String(fmt.Sprintf("%f", data.RevenuePerEmployee)),
			},
			"DebtToEquity": {
				N: aws.String(fmt.Sprintf("%f", data.DebtToEquity)),
			},
			"ProfitMargin": {
				N: aws.String(fmt.Sprintf("%f", data.ProfitMargin)),
			},
			"EnterpriseValue": {
				N: aws.String(fmt.Sprintf("%f", data.EnterpriseValue)),
			},
			"EnterpriseValueToRevenue": {
				N: aws.String(fmt.Sprintf("%f", data.EnterpriseValueToRevenue)),
			},
			"PriceToSales": {
				N: aws.String(fmt.Sprintf("%f", data.PriceToSales)),
			},
			"PriceToBook": {
				N: aws.String(fmt.Sprintf("%f", data.PriceToBook)),
			},
			"ForwardPERatio": {
				N: aws.String(fmt.Sprintf("%f", data.ForwardPERatio)),
			},
			"PEGRatio": {
				N: aws.String(fmt.Sprintf("%f", data.PEGRatio)),
			},
			"PEHigh": {
				N: aws.String(fmt.Sprintf("%f", data.PEHigh)),
			},
			"PELow": {
				N: aws.String(fmt.Sprintf("%f", data.PELow)),
			},
			"Week52HighDate": {
				S: aws.String(time.Time(data.Week52HighDate).Format("2006-01-02")),
			},
			"Week52LowDate": {
				S: aws.String(time.Time(data.Week52LowDate).Format("2006-01-02")),
			},
			"PutCallRatio": {
				N: aws.String(fmt.Sprintf("%f", data.PutCallRatio)),
			},
		},
	}
	if _, err := ddbClient.PutItem(input); err != nil {
		return err
	}
	log.Infof("Saved stats for %s", symbol)
	return nil
}

func handler(e events.DynamoDBEvent) error {
	// Loop through new records acting only on insert
	var item map[string]events.DynamoDBAttributeValue
	var tableName string
	for _, v := range e.Records {
		switch v.EventName {
		case "INSERT", "MODIFY":
			tableName = strings.Split(v.EventSourceArn, "/")[1]
			log.WithFields(logrus.Fields{"EventID": v.EventID}).Infof("Processing a %s from %s table", v.EventName, tableName)
			item = v.Change.NewImage
			if err := processItem(item); err != nil {
				return err
			}
			log.WithFields(logrus.Fields{"EventID": v.EventID}).Infof("Finished processing a %s from %s table", v.EventName, tableName)
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
