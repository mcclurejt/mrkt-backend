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
	writeRequests := make([]*ddb.WriteRequest, len(historical))
	for i, data := range historical {
		writeRequests[i] = &ddb.WriteRequest{
			PutRequest: &ddb.PutRequest{
				Item: map[string]*ddb.AttributeValue{
					"Date": {
						S: aws.String(data.Date),
					},
					"Symbol": {
						S: aws.String(symbol.String()),
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
						N: aws.String(fmt.Sprintf("%f", data.Close-data.Open)),
					},
					"ChangePercent": {
						N: aws.String(fmt.Sprintf("%f", (data.Close-data.Open)/data.Open)),
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
	var eventID string
	for _, v := range e.Records {
		switch v.EventName {
		case "INSERT", "UPDATE":
			tableName = strings.Split(v.EventSourceArn, "/")[1]
			eventID = v.EventID
			log.WithFields(logrus.Fields{"EventID": eventID}).Infof("Processing an insert/update from %s table", tableName)
			item = v.Change.NewImage
			if err := processItem(item); err != nil {
				return err
			}
		}
	}
	log.WithFields(logrus.Fields{"EventID": eventID}).Infof("Finished processing an insert/update from %s table", tableName)
	return nil
}

func main() {
	lambda.Start(handler)
}
