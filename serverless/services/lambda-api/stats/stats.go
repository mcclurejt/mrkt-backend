package main

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	iex "github.com/goinvest/iexcloud/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

type StatsWithSymbol struct {
	iex.AdvancedStats
	Symbol           string
	NextDividendDate string `json:"NextDividendDate"`
	ExDividendDate   string `json:"ExDividendDate"`
	NextEarningsDate string `json:"NextEarningsDate"`
	Week52HighDate   string `json:"Week52HighDate"`
	Week52LowDate    string `json:"Week52LowDate"`
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
	log.Infof("Retrieving Stats for %s...", symbol)
	out, err := ddbClient.GetItem(
		&ddb.GetItemInput{
			TableName: aws.String("Stats"),
			Key: map[string]*ddb.AttributeValue{
				"Symbol": {S: aws.String(symbol)},
			},
		})
	// parse the response object
	stats := StatsWithSymbol{}
	err = dynamodbattribute.UnmarshalMap(out.Item, &stats)
	if err != nil {
		return util.ErrorToGatewayResponse(err)
	}
	log.Infof("Retrieved stats for symbol %s", symbol)
	return util.ObjectToGatewayResponse(stats)
}

func main() {
	lambda.Start(handler)
}
