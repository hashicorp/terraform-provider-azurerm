package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTableDestination struct {
	Name                     *string `json:"name,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
	TableName                *string `json:"tableName,omitempty"`
}
