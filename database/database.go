package database

// AssetPropertyService - Interface defining a service that provides unique information of an asset. These properties generally reside in their own table and are linked to the assed by Ticker ID
type AssetPropertyService interface {
	GetTableName() string
	GetTableColumns() []string
}

// SQLClient - Interface defining the functions of a generic SQL Client, used to interact with the db
type SQLClient interface {
	HasTable(tableName string) bool
	CreateTable(aps AssetPropertyService) error
	CreateAllTables(client interface{}) error
	DropTable(aps AssetPropertyService) error
	Insert(tableName string, headers []string, values []interface{}) error
	GetTickerID(ticker string) (int, error)
}
