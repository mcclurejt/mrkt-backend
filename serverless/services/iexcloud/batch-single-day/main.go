package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/mcclurejt/mrkt-backend/api/iexcloud"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var (
	iexClient *iexcloud.IEXCloudClient
	ddbClient dynamodbiface.DynamoDBAPI
	wg        sync.WaitGroup
)

func init() {
	iexClient = iexcloud.NewIEXCloudClient("pk_1d8a2228abd84b0598a6cf91a5d09f63")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		return
	}
	ddbClient = ddb.New(awsSession)
}

func apiResponse(status int, body interface{}) (*Response, error) {
	resp := Response{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp.StatusCode = status
	strBody, _ := json.Marshal(body)
	resp.Body = string(strBody)
	return &resp, nil
}

func PostChartSingleDay(s iexcloud.OHLCV) error {
	defer wg.Done()
	av, err := dynamodbattribute.MarshalMap(s)
	if err != nil {
		return err
	}
	input := &ddb.PutItemInput{
		Item:      av,
		TableName: aws.String("ChartSingleDay"),
	}
	_, err = ddbClient.PutItem(input)
	if err != nil {
		fmt.Println("D")
		return err
	}
	return nil
}

func GetPutRequests(charts []iexcloud.OHLCV) ([]*ddb.PutRequest, error) {
	requests := []*ddb.PutRequest{}
	for _, symbol := range charts {
		av, err := dynamodbattribute.MarshalMap(symbol)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &ddb.PutRequest{Item: av})
	}
	return requests, nil
}

func GetBatchWriteRequests(requests []*ddb.PutRequest, tableName string) []*ddb.BatchWriteItemInput {
	inputs := []*ddb.BatchWriteItemInput{}
	l := len(requests)
	batchSize := 25
	var batch []*ddb.PutRequest
	for i := 0; i < l; i += batchSize {
		if i+batchSize > l {
			batch = requests[i:l]
		} else {
			batch = requests[i : i+batchSize]
		}
		requestItems := map[string][]*ddb.WriteRequest{}
		requestItems[tableName] = []*ddb.WriteRequest{}
		for _, item := range batch {
			requestItems[tableName] = append(requestItems[tableName], &ddb.WriteRequest{PutRequest: item})
		}
		inputs = append(inputs, &ddb.BatchWriteItemInput{
			RequestItems: requestItems,
		})
	}
	return inputs
}

func ExecuteBatchWriteRequest(req *ddb.BatchWriteItemInput) {
	defer wg.Done()
	_, err := ddbClient.BatchWriteItem(req)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (*Response, error) {
	symbols, err := iexClient.IexSymbols.Get(context.Background())
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}

	symbolList := []string{}
	for i := 0; i < 50; i++ {
		symbolList = append(symbolList, symbols[i].Symbol)
	}
	ohlcvs, err := iexClient.Chart.GetBatchSingleDay(context.Background(), symbolList, "20201103")
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	fmt.Println(ohlcvs)

	// make the putrequests
	putRequests, err := GetPutRequests(ohlcvs)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	// make the batchwriterequests
	batchWriteRequests := GetBatchWriteRequests(putRequests, "ChartSingleDay")
	// post to ddb
	for _, v := range batchWriteRequests {
		wg.Add(1)
		go ExecuteBatchWriteRequest(v)
	}
	wg.Wait()

	// fmt.Println("Finished fetching and posting to ddb")
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"success": fmt.Sprintf("Successfully updated %d records in ChartSingleDay \n", len(putRequests)),
	})
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	json.HTMLEscape(&buf, body)

	return apiResponse(http.StatusOK, buf.String())
}

func main() {
	lambda.Start(Handler)
}
