package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	av "github.com/mcclurejt/mrkt-backend/api/alphavantage"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var avClient av.AlphaVantageClient

func init() {
	avClient = av.NewAlphaVantageClient("LXCN06KPP1KPOYC2")
}

type Input struct {
	Ticker string
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, input Input) (Response, error) {
	var buf bytes.Buffer
	ticker := input.Ticker

	ts, err := avClient.MonthlyAdjustedTimeSeries.Get(&av.MonthlyAdjustedTimeSeriesOptions{Symbol: ticker})
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	body, err := json.Marshal(ts)
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}