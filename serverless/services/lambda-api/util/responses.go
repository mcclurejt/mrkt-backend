package util

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func EncodeStringAsBody(s string) string {
	body, err := json.Marshal(map[string]string{"message": s})
	if err != nil {
		return "{ \"message\": \"failed marshaling message text\" }"
	}
	return string(body)
}

func ErrorToGatewayResponse(err error) (events.APIGatewayProxyResponse, error) {
	aerr, ok := err.(awserr.Error)
	// Return if not aws error
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Body:       fmt.Sprintf("An unknown error occurred: %s", err.Error()),
		}, err
	}
	// Return gateway response with details filled by error
	return events.APIGatewayProxyResponse{
		StatusCode: 501,
		Body:       fmt.Sprintf("{ \"message\": \"%s\" }", aerr.Message()),
	}, err
}

func ObjectToGatewayResponse(v interface{}) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return ErrorToGatewayResponse(err)
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}
	return response, nil
}
