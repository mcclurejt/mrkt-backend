package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	iex "github.com/goinvest/iexcloud/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mcclurejt/mrkt-backend/serverless/services/lambda-api/util"
	"github.com/sirupsen/logrus"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type HistoricalWithSymbol struct {
	iex.HistoricalDataPoint
	Symbol string
}

var (
	ddbClient dynamodbiface.DynamoDBAPI
	log       *logrus.Logger
)

func init() {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return
	}
	ddbClient = ddb.New(awsSession)
	log = logrus.New()
}

func historicalForSymbol(symbol string) ([]HistoricalWithSymbol, error) {
	historical := []HistoricalWithSymbol{}
	err := ddbClient.QueryPages(&ddb.QueryInput{
		TableName:              aws.String("Historical"),
		KeyConditionExpression: aws.String("#pk = :s"),
		ExpressionAttributeNames: map[string]*string{
			"#pk": aws.String("Symbol"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(symbol),
			},
		},
	}, func(page *dynamodb.QueryOutput, _ bool) bool {
		h := []HistoricalWithSymbol{}
		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &h)
		if err != nil {
			log.Errorf("Unable to unmarshal AWS data: err = %v", err)
			return true
		}
		historical = append(historical, h...)
		return true
	})
	if err != nil {
		return historical, err
	}
	if len(historical) == 0 {
		return historical, util.NewErrorDataNotFoundForSymbol("Historical", symbol)
	}
	return historical, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log := log.WithFields(logrus.Fields{"path": request.Path, "method": request.HTTPMethod})
	if request.HTTPMethod != http.MethodGet {
		err := util.NewErrorMethodNotImplemented(request.HTTPMethod)
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Body:       util.EncodeStringAsBody(err.Error()),
		}, err
	}
	// extract symbol from the path
	symbol := request.PathParameters["symbol"]
	symbol = strings.ToUpper(symbol)
	// get historical data for symbol
	log.Infof("Retrieving Historical Data for %s...", symbol)
	historical, err := historicalForSymbol(symbol)
	if err != nil {
		return util.ErrorToGatewayResponse(err)
	}
	log.Infof("Retrieved %d datapoints for symbol %s", len(historical), symbol)
	return util.ObjectToGatewayResponse(historical)
}

func main() {
	lambda.Start(handler)
}
