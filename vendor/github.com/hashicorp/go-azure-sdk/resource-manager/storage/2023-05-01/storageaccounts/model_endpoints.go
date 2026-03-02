package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Endpoints struct {
	Blob               *string                           `json:"blob,omitempty"`
	Dfs                *string                           `json:"dfs,omitempty"`
	File               *string                           `json:"file,omitempty"`
	InternetEndpoints  *StorageAccountInternetEndpoints  `json:"internetEndpoints,omitempty"`
	MicrosoftEndpoints *StorageAccountMicrosoftEndpoints `json:"microsoftEndpoints,omitempty"`
	Queue              *string                           `json:"queue,omitempty"`
	Table              *string                           `json:"table,omitempty"`
	Web                *string                           `json:"web,omitempty"`
}
