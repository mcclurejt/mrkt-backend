package coin

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	gn "github.com/mcclurejt/mrkt-backend/api/glassnode"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var gnClient *gn.GlassNodeClient

func init() {
	gnClient = gn.NewGlassNodeClient("105d32cc-afc0-4358-b335-891a35e80736")
}

func Handler(ctx context.Context, options *gn.CoinOptions) (Response, error) {
	nupl, err := gnClient.NetUnrealizedProfitLoss.Get(options)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var buf bytes.Buffer
	body, err := json.Marshal(nupl)
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
