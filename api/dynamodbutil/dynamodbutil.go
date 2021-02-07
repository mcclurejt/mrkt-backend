package dynamodbutil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const DefaultBillingMode = db.BillingModePayPerRequest
const MaxBatchSize = 25

var ValidAttributeTypeMap = map[string]bool{
	"S": true,
	"N": true,
	"B": true,
}

var ValidKeyTypeMap = map[string]bool{
	"HASH":  true,
	"RANGE": true,
}

var AttributeNameTags = []string{"attributename", "an"}

var AttributeTypeTags = []string{"attributetype", "at"}

var KeyTypeTags = []string{"keytype", "kt"}

func CreateHistoricalTable(ddbClient dynamodbiface.DynamoDBAPI) (*db.CreateTableOutput, error) {
	input := &dynamodb.CreateTableInput{
		TableName:   aws.String("Historical"),
		BillingMode: aws.String(DefaultBillingMode),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Symbol"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Open"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("High"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Low"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Close"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Volume"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("Change"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("ChangePercent"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Date"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Symbol"),
				KeyType:       aws.String("RANGE"),
			},
		},
	}
	return ddbClient.CreateTable(input)
}

// CreateTableInputFromStruct - Generates CreateTableInput using the tags contained in the input struct s, if s is not a struct, an error is thrown
// AttributeName: Uses the field's name
// AttributeType: Accepts both `at` and `attributetype` struct tags.
// KeyType: Accepts both `kt` and `keytype` struct tags
func CreateTableInputFromStruct(s interface{}) (*dynamodb.CreateTableInput, error) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("Error: Input must be a struct or pointer to a struct")
	}
	tableName := t.Name()
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{},
		KeySchema:            []*dynamodb.KeySchemaElement{},
		BillingMode:          aws.String(DefaultBillingMode),
		TableName:            aws.String(tableName),
	}
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		// Attribute Name
		attributeName := f.Name
		// Attribute Type
		tagFound := false
		attributeType := ""
		for _, tag := range AttributeTypeTags {
			val, ok := f.Tag.Lookup(tag)
			if ok {
				attributeType = val
				tagFound = true
			}
		}
		if !tagFound {
			continue
		}
		if _, ok := ValidAttributeTypeMap[attributeType]; !ok {
			return nil, fmt.Errorf("Error: %s is not a valid Attribute Type", attributeType)
		}
		// Key Type
		tagFound = false
		keyType := ""
		for _, tag := range KeyTypeTags {
			val, ok := f.Tag.Lookup(tag)
			if ok {
				keyType = val
				tagFound = true
			}
		}
		if !tagFound {
			continue
		}
		if _, ok := ValidKeyTypeMap[keyType]; !ok {
			return nil, fmt.Errorf("Error: %s is not a valid Key Type", keyType)
		}
		// Add items to CreateTableInput
		input.AttributeDefinitions = append(input.AttributeDefinitions, &dynamodb.AttributeDefinition{
			AttributeName: aws.String(attributeName),
			AttributeType: aws.String(attributeType),
		})
		input.KeySchema = append(input.KeySchema, &dynamodb.KeySchemaElement{
			AttributeName: aws.String(attributeName),
			KeyType:       aws.String(keyType),
		})
	}
	return input, nil
}

// PutItemInputFromStruct - Generates PutItemInput from the given struct
func PutItemInputFromStruct(item interface{}) (*dynamodb.PutItemInput, error) {
	// create the attributevalue
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, err
	}
	// check if it contains any inferred attribute type tags. If so, modify that field accordingly
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.Name()),
	}
	return input, nil
}

// PutRequestsFromSlice - Converts a slice of structs to an array of PutRequests
func PutRequestsFromSlice(arr interface{}) ([]*db.PutRequest, error) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("Error: Input argument must be a slice")
	}
	requests := []*db.PutRequest{}
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		av, err := dynamodbattribute.MarshalMap(item.Interface())
		if err != nil {
			return nil, err
		}
		requests = append(requests, &db.PutRequest{Item: av})
	}
	return requests, nil
}

// ConvertToBatchPutRequest - Consolidates the input slice of PutRequests to a BatchWriteRequest
func ConvertToBatchPutRequest(requests []*db.PutRequest, tableName string) []*db.BatchWriteItemInput {
	inputs := []*db.BatchWriteItemInput{}
	l := len(requests)
	var batch []*db.PutRequest
	for i := 0; i < l; i += MaxBatchSize {
		if i+MaxBatchSize > l {
			batch = requests[i:l]
		} else {
			batch = requests[i : i+MaxBatchSize]
		}
		requestItems := map[string][]*db.WriteRequest{}
		requestItems[tableName] = []*db.WriteRequest{}
		for _, item := range batch {
			requestItems[tableName] = append(requestItems[tableName], &db.WriteRequest{PutRequest: item})
		}
		inputs = append(inputs, &db.BatchWriteItemInput{
			RequestItems: requestItems,
		})
	}
	return inputs
}
