package outputs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureTableOutputDataSourceProperties struct {
	AccountKey      *string   `json:"accountKey,omitempty"`
	AccountName     *string   `json:"accountName,omitempty"`
	BatchSize       *int64    `json:"batchSize,omitempty"`
	ColumnsToRemove *[]string `json:"columnsToRemove,omitempty"`
	PartitionKey    *string   `json:"partitionKey,omitempty"`
	RowKey          *string   `json:"rowKey,omitempty"`
	Table           *string   `json:"table,omitempty"`
}
