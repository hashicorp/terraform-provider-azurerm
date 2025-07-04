package datasets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMySqlTableDatasetTypeProperties struct {
	Table     *interface{} `json:"table,omitempty"`
	TableName *interface{} `json:"tableName,omitempty"`
}
