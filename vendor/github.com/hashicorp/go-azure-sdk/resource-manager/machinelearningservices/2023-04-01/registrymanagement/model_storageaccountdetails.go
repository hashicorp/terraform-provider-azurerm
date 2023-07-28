package registrymanagement

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountDetails struct {
	SystemCreatedStorageAccount *SystemCreatedStorageAccount `json:"systemCreatedStorageAccount,omitempty"`
	UserCreatedStorageAccount   *UserCreatedStorageAccount   `json:"userCreatedStorageAccount,omitempty"`
}
