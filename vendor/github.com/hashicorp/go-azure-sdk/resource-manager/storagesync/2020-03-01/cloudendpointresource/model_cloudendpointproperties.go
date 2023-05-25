package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudEndpointProperties struct {
	AzureFileShareName       *string `json:"azureFileShareName,omitempty"`
	BackupEnabled            *string `json:"backupEnabled,omitempty"`
	FriendlyName             *string `json:"friendlyName,omitempty"`
	LastOperationName        *string `json:"lastOperationName,omitempty"`
	LastWorkflowId           *string `json:"lastWorkflowId,omitempty"`
	PartnershipId            *string `json:"partnershipId,omitempty"`
	ProvisioningState        *string `json:"provisioningState,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
	StorageAccountTenantId   *string `json:"storageAccountTenantId,omitempty"`
}
