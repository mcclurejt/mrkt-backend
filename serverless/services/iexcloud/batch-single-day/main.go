package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ddb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	util "github.com/mcclurejt/mrkt-backend/api/dynamodbutil"
	"github.com/mcclurejt/mrkt-backend/api/iexcloud"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

const TableName = "ChartSingleDay"

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

func ExecuteBatchWriteRequest(req *ddb.BatchWriteItemInput, errCh chan error) {
	defer wg.Done()
	_, err := ddbClient.BatchWriteItem(req)
	if err != nil {
		errCh <- err
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
	t1 := time.Now()
	symbolList := []string{}
	for i := 0; i < len(symbols); i++ {
		symbolList = append(symbolList, symbols[i].Symbol)
	}
	t2 := time.Now()
	ohlcvs, err := iexClient.Chart.GetBatchSingleDay(context.Background(), symbolList, "20201103")
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	t3 := time.Now()
	// make the putrequests
	putRequests, err := util.PutRequestsFromSlice(ohlcvs)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	// make the batchwriterequests
	batchWriteRequests := util.ConvertToBatchPutRequest(putRequests, TableName)
	// post to ddb
	errCh := make(chan error)
	for _, v := range batchWriteRequests {
		wg.Add(1)
		go ExecuteBatchWriteRequest(v, errCh)
	}
	// close out once waitgroup is done
	wg.Wait()
	close(errCh)

	fmt.Printf("symbols: %d", t2.Sub(t1).Seconds())
	fmt.Printf("ohlcv: %d", t3.Sub(t2).Seconds())

	// collect all errors from the goroutines and return
	if len(errCh) > 0 {
		errString := ""
		for e := range errCh {
			errString += e.Error() + ", "
		}
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(errString)})
	}

	// lambda was successful
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
