package dynamodb

import (
	"errors"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	db "github.com/aws/aws-sdk-go/service/dynamodb"
	attribute "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const DefaultBillingMode = db.BillingModePayPerRequest

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
