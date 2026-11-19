// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

type MssqlDatabaseSensitivityLabelModel struct {
	ColumnName        string  `tfschema:"column_name"`
	DatabaseId        string  `tfschema:"database_id"`
	SchemaName        string  `tfschema:"schema_name"`
	TableName         string  `tfschema:"table_name"`
	InformationType   *string `tfschema:"information_type"`
	InformationTypeId *string `tfschema:"information_type_id"`
	LabelId           *string `tfschema:"label_id"`
	LabelName         *string `tfschema:"label_name"`
	Rank              *string `tfschema:"rank"`
}

func (Resource) ModelObject() any {
	return &MssqlDatabaseSensitivityLabelModel{}
}
