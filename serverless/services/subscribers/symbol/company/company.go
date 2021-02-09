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
	// Retrieve historical data
	symbol, ok := item["Symbol"]
	if !ok {
		return errors.New("Symbol Key Not Found")
	}
	log.Info("Retrieving company summary for %s", symbol.String())
	t := time.Now()
	data, err := iexClient.Company(context.Background(), symbol.String())
	if err != nil {
		return err
	}
	log.Infof("Retrieved company summary in %.2fs", time.Now().Sub(t).Seconds())
	// Form the list of requests
	input := &ddb.PutItemInput{
		TableName: aws.String("Company"),
		Item: map[string]*ddb.AttributeValue{
			"Symbol": {
				S: aws.String(symbol.String()),
			},
			"CompanyName": {
				S: aws.String(data.Name),
			},
			"Exchange": {
				S: aws.String(data.Exchange),
			},
			"Industry": {
				S: aws.String(data.Industry),
			},
			"Website": {
				S: aws.String(data.Website),
			},
			"Description": {
				S: aws.String(data.Description),
			},
			"CEO": {
				S: aws.String(data.CEO),
			},
			"SecurityName": {
				S: aws.String(data.SecurityName),
			},
			"IssueType": {
				S: aws.String(data.IssueType),
			},
			"Sector": {
				S: aws.String(data.Sector),
			},
			"PrimarySicCode": {
				N: aws.String(fmt.Sprintf("%d", data.PrimarySICCode)),
			},
			"Employees": {
				N: aws.String(fmt.Sprintf("%d", data.Employees)),
			},
			"Tags": {
				SS: aws.StringSlice(data.Tags),
			},
			"Address": {
				S: aws.String(data.Address),
			},
			"Address2": {
				S: aws.String(data.Address2),
			},
			"State": {
				S: aws.String(data.State),
			},
			"City": {
				S: aws.String(data.City),
			},
			"Zip": {
				S: aws.String(data.Zip),
			},
			"Country": {
				S: aws.String(data.Country),
			},
			"Phone": {
				S: aws.String(data.Phone),
			},
		},
	}
	if _, err := ddbClient.PutItem(input); err != nil {
		return err
	}
	log.Infof("Saved company summary for %s", symbol)
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
