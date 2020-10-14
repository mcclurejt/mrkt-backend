package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	av "github.com/mcclurejt/mrkt-backend/api/alphavantage"
	"github.com/mcclurejt/mrkt-backend/api/dynamodb"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
type Input struct {
	Body string `json:"body"`
}

var avClient av.AlphaVantageClient
var ddbClient *dynamodb.Client

func init() {
	avClient = av.NewAlphaVantageClient("LXCN06KPP1KPOYC2")
	ddbClient = dynamodb.New()
}

func Handler(ctx context.Context, input Input) (Response, error) {
	var ts = av.MonthlyAdjustedTimeSeries{}
	err := json.Unmarshal([]byte(input.Body), &ts)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	entries := ts.TimeSeries
	for _, v := range entries {
		err = ddbClient.PutItem(avClient.MonthlyAdjustedTimeSeries, v)
		if err != nil {
			return Response{StatusCode: 404}, err
		}
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintln("Success! Uploaded to dynamodb"),
	}
	fmt.Println(resp)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
