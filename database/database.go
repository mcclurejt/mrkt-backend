package database

type Client interface {
	Insert(tableName string, headers []string, values []interface{}) error
}
