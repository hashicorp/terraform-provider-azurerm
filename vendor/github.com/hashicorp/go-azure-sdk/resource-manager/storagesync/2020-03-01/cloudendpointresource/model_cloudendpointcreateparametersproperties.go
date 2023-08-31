package cloudendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudEndpointCreateParametersProperties struct {
	AzureFileShareName       *string `json:"azureFileShareName,omitempty"`
	FriendlyName             *string `json:"friendlyName,omitempty"`
	StorageAccountResourceId *string `json:"storageAccountResourceId,omitempty"`
	StorageAccountTenantId   *string `json:"storageAccountTenantId,omitempty"`
}
