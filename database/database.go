package database

type Client interface {
	Insert(tableName string, headers []string, values []interface{}) error
	CreateTable(tableName string, columns []string) error
	DropTable(tableName string) error
}
