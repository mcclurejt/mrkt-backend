package main

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	"github.com/mcclurejt/mrkt-backend/serverless/services/lambda-api/util"
	"github.com/sirupsen/logrus"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

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

func listSymbols() ([]string, error) {
	symbols := []string{}
	for {
		symbolBatch, err := ddbClient.Scan(&ddb.ScanInput{TableName: aws.String("Symbols")})
		if err != nil {
			return []string{}, err
		}
		for _, attributeMap := range symbolBatch.Items {
			log.Infof("%v", attributeMap)
			if symbol, ok := attributeMap["Symbol"]; ok && symbol.S != nil {
				symbols = append(symbols, *symbol.S)
			}
		}
		if len(symbolBatch.LastEvaluatedKey) == 0 {
			break
		}
	}
	return symbols, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log := log.WithFields(logrus.Fields{"path": request.Path, "method": request.HTTPMethod})
	// check http method
	if request.HTTPMethod != http.MethodGet {
		errorText := "request method not implemented for route"
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Body:       util.EncodeStringAsBody(errorText),
		}, errors.New(errorText)
	}
	// get list of symbols
	log.Info("Retrieving Symbols...")
	symbols, err := listSymbols()
	if err != nil {
		return util.ErrorToGatewayResponse(err)
	}
	log.Infof("Retrieved %d symbols", len(symbols))
	return util.ObjectToGatewayResponse(symbols)
}

func main() {
	lambda.Start(handler)
}
