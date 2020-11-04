package dynamodb

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	db "github.com/aws/aws-sdk-go/service/dynamodb"
	attribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const DefaultBillingMode = db.BillingModePayPerRequest

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

type DataProvider interface {
	GetCreateTableInput() *db.CreateTableInput
	GetPutItemInput() *db.PutItemInput
}

type Client struct {
	base *db.DynamoDB
}

// New - Creates a new instance of the dynamodb client
func New() *Client {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	return &Client{base: db.New(sess)}
}

// CreateTable - Creates a new table from the input specified by the DynamodbDataProvider
func (c *Client) CreateTable(p DataProvider) error {
	input := p.GetCreateTableInput()
	_, err := c.base.CreateTable(input)
	if err != nil {
		return err
	}
	return nil
}

// CreateAllTables - Creates a new table for each DynamodbDataProvider contained by the client
func (c *Client) CreateAllTables(client interface{}) error {
	v := reflect.ValueOf(client)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.CanInterface() {
			p, ok := f.Interface().(DataProvider)
			if ok {
				err := c.CreateTable(p)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ListTables - Lists all tables present in dynamodb
func (c *Client) ListTables() (*db.ListTablesOutput, error) {
	input := &db.ListTablesInput{}
	return c.base.ListTables(input)
}

// PutItem - Adds an item in the specified table using the configuration provided by the DataProvider's GetPutItemInput()
func (c *Client) PutItem(p DataProvider, item interface{}) error {
	av, err := attribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := p.GetPutItemInput()
	input.Item = av
	_, err = c.base.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

// PutAllItems - Adds each item from the list using the configuration provided by the DataProvider
func (c *Client) PutAllItems(p DataProvider, items interface{}) error {
	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		return errors.New("Must pass in a slice")
	}

	slice := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		slice[i] = val.Index(i).Interface()
	}

	var err error
	for _, item := range slice {
		err = c.PutItem(p, item)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateTableInputFromStruct - Generates CreateTableInput using the tags contained in the input struct s, if s is not a struct, an error is thrown
// AttributeName: Accepts both `an` and `attributename` struct tags
// AttributeType: Accepts both `at` and `attributetype` struct tags
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
		tagFound := false
		attributeName := ""
		for _, tag := range AttributeNameTags {
			val, ok := f.Tag.Lookup(tag)
			if ok {
				attributeName = val
				tagFound = true
			}
		}
		if !tagFound {
			continue
		}
		// Attribute Type
		tagFound = false
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
