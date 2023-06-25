package dataset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlDBTableProperties struct {
	DataSetId           *string `json:"dataSetId,omitempty"`
	DatabaseName        string  `json:"databaseName"`
	SchemaName          string  `json:"schemaName"`
	SqlServerResourceId string  `json:"sqlServerResourceId"`
	TableName           string  `json:"tableName"`
}
