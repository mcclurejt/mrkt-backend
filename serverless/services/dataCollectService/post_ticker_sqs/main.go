package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mcclurejt/mrkt-backend/api"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var msClient api.MarketStackClient

func init() {
	msClient = api.NewMarketStackClient("02378e09665e4a13b514d5cb29855994")
}

func GetQueueUrl(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	fmt.Println("Getting Queue URL")
	svc := sqs.New(sess)
	url, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	fmt.Printf("Queue URL: %s", *url.QueueUrl)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func SendMsg(sess *session.Session, input *sqs.SendMessageBatchInput) error {
	// Create an SQS service client
	svc := sqs.New(sess)
	_, err := svc.SendMessageBatch(input)
	if err != nil {
		return err
	}
	return nil
}

type Input struct {
	Body []string `json:"Body"` // array of Tickers
}

func Handler(ctx context.Context, input Input) (Response, error) {
	sess := session.New()

	url := fmt.Sprintf("ticker-data-collection")
	queueUrl, err := GetQueueUrl(sess, &url)
	if err != nil {
		return Response{}, err
	}

	entries := make([]*sqs.SendMessageBatchRequestEntry, len(input.Body))
	for i, ticker := range input.Body {
		entries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:           aws.String(ticker),
			DelaySeconds: aws.Int64(10),
			MessageBody:  aws.String(ticker),
		}
	}

	sqsInput := &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: queueUrl.QueueUrl,
	}
	err = SendMsg(sess, sqsInput)
	if err != nil {
		return Response{}, err
	}
	fmt.Println("Sent tickers to queue")

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
