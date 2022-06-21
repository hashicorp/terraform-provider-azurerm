package linkedstorageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedStorageAccountsProperties struct {
	DataSourceType    *DataSourceType `json:"dataSourceType,omitempty"`
	StorageAccountIds *[]string       `json:"storageAccountIds,omitempty"`
}
