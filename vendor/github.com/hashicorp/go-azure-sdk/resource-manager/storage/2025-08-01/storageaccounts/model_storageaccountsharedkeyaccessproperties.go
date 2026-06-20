package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountSharedKeyAccessProperties struct {
	Blob  *ServiceSharedKeyAccessProperties `json:"blob,omitempty"`
	File  *ServiceSharedKeyAccessProperties `json:"file,omitempty"`
	Queue *ServiceSharedKeyAccessProperties `json:"queue,omitempty"`
	Table *ServiceSharedKeyAccessProperties `json:"table,omitempty"`
}
