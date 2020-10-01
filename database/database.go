package database

type SQLClient interface {
	Insert(tableName string, headers []string, values []interface{}) error
	HasTable(tableName string) bool
	CreateTable(tableName string, columns []string) error
	DropTable(tableName string) error
	GetTickerID(ticker string) (int, error)
}
