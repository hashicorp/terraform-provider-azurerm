package linkedstorageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedStorageAccountsResource struct {
	Id         *string                         `json:"id,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Properties LinkedStorageAccountsProperties `json:"properties"`
	Type       *string                         `json:"type,omitempty"`
}
