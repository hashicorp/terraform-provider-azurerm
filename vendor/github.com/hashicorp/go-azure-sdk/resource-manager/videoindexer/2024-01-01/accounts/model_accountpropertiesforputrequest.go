package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountPropertiesForPutRequest struct {
	AccountId           *string                       `json:"accountId,omitempty"`
	AccountName         *string                       `json:"accountName,omitempty"`
	ProvisioningState   *ProvisioningState            `json:"provisioningState,omitempty"`
	StorageServices     *StorageServicesForPutRequest `json:"storageServices,omitempty"`
	TenantId            *string                       `json:"tenantId,omitempty"`
	TotalSecondsIndexed *int64                        `json:"totalSecondsIndexed,omitempty"`
}
