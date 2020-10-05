package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mcclurejt/mrkt-backend/api"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var avClient api.AlphaVantageClient

func init() {
	avClient = api.NewAlphaVantageClient("LXCN06KPP1KPOYC2")
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	ts, err := avClient.MonthlyAdjustedTimeSeriesService.Get("BRK-A")
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
	fmt.Println(resp)
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}