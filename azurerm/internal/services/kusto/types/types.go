package dataplane

import "github.com/Azure/azure-kusto-go/kusto/data/value"

type KustoTableRecord struct {
	TableName    string
	DatabaseName string
	Folder       string
	DocString    string
}

type KustoTableColumnSchemaRecord struct {
	ColumnName    string
	ColumnOrdinal value.Int
	DataType      string
	ColumnType    string
}
