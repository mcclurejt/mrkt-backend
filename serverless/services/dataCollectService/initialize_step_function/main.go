package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type StateMachineInput struct {
	Ticker string
}

func InvokeStepFunction(sess *session.Session, input *sfn.StartExecutionInput) (*sfn.StartExecutionOutput, error) {
	svc := sfn.New(sess, &aws.Config{})
	req, resp := svc.StartExecutionRequest(input)
	err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event events.SQSEvent) (Response, error) {
	sess := session.New()
	if len(event.Records) == 0 {
		return Response{}, errors.New("No SQS message passed to function")
	}
	t := StateMachineInput{
		Ticker: event.Records[0].Body,
	}
	var jsonData []byte
	jsonData, err := json.Marshal(t)
	if err != nil {
		return Response{}, errors.New("Error marshalling ticker into json")
	}

	input := &sfn.StartExecutionInput{
		StateMachineArn: aws.String("arn:aws:states:us-west-2:115333527451:stateMachine:DataCollectionStepFunctionsStateMachine-0Hs5igHvGMak"),
		Input:           aws.String(string(jsonData)),
	}

	_, err = InvokeStepFunction(sess, input)
	if err != nil {
		return Response{}, errors.New("Error invoking step function")
	}

	resp := Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
