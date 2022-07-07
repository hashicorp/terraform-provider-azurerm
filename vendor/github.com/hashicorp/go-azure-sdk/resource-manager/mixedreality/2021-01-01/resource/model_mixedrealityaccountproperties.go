package resource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MixedRealityAccountProperties struct {
	AccountDomain      *string `json:"accountDomain,omitempty"`
	AccountId          *string `json:"accountId,omitempty"`
	StorageAccountName *string `json:"storageAccountName,omitempty"`
}
