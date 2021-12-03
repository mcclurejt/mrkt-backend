package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mcclurejt/mrkt-backend/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	iex "github.com/goinvest/iexcloud/v2"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
	conf := config.New() //env
	iexClient = iex.NewClient(conf.Api.IEXCloudAPIKey)
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
	// Retrieve historical data
	symbol, ok := item["Symbol"]
	if !ok {
		return errors.New("Symbol Key Not Found")
	}
	log.Info("Retrieving historical data for %s", symbol.String())
	t := time.Now()
	historical, err := iexClient.HistoricalPrices(context.Background(), symbol.String(), iex.SixMonthHistorical, &iex.HistoricalOptions{ChangeFromClose: true})
	if err != nil {
		return err
	}
	log.Infof("Retrieved %d historical datapoints in %.2fs", len(historical), time.Now().Sub(t).Seconds())
	// Form the list of requests
	var change, changePercent float64
	writeRequests := make([]*ddb.WriteRequest, len(historical))
	for i, data := range historical {
		if i == 0 {
			change = data.Close - data.Open
			changePercent = change / data.Open
		} else {
			change = data.Close - historical[i-1].Close
			changePercent = change / historical[i-1].Close
		}
		writeRequests[i] = &ddb.WriteRequest{
			PutRequest: &ddb.PutRequest{
				Item: map[string]*ddb.AttributeValue{
					"Symbol": {
						S: aws.String(symbol.String()),
					},
					"Date": {
						S: aws.String(data.Date),
					},
					"Open": {
						N: aws.String(fmt.Sprintf("%f", data.Open)),
					},
					"High": {
						N: aws.String(fmt.Sprintf("%f", data.High)),
					},
					"Low": {
						N: aws.String(fmt.Sprintf("%f", data.Low)),
					},
					"Close": {
						N: aws.String(fmt.Sprintf("%f", data.Close)),
					},
					"Volume": {
						N: aws.String(fmt.Sprintf("%d", data.Volume)),
					},
					"Change": {
						N: aws.String(fmt.Sprintf("%f", change)),
					},
					"ChangePercent": {
						N: aws.String(fmt.Sprintf("%f", changePercent)),
					},
				},
			},
		}
	}
	// Launch goroutines to execute requests in batches of 25
	errs, _ := errgroup.WithContext(context.Background())
	for i := 0; i < len(historical); i += 25 {
		j := i + 25
		if j > len(historical) {
			j = len(historical)
		}
		reqs := writeRequests[i:j]
		errs.Go(func() error {
			return executeBatch(reqs)
		})
	}
	// Wait until all requests are done and return
	return errs.Wait()
}

func executeBatch(reqs []*ddb.WriteRequest) error {
	t := time.Now()
	batchRequest := &ddb.BatchWriteItemInput{
		RequestItems: map[string][]*ddb.WriteRequest{
			"Historical": reqs,
		},
	}
	if _, err := ddbClient.BatchWriteItem(batchRequest); err != nil {
		return err
	}
	log.Infof("Executed batch request in %.2fs", time.Now().Sub(t).Seconds())
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
