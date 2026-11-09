package connectionanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountSettings struct {
	Path             *string `json:"path,omitempty"`
	StorageAccountId string  `json:"storageAccountId"`
}
