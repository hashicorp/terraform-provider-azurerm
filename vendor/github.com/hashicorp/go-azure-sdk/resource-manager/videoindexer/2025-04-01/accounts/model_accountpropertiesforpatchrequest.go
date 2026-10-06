package accounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountPropertiesForPatchRequest struct {
	AccountId                  *string                         `json:"accountId,omitempty"`
	OpenAiServices             *OpenAiServicesForPatchRequest  `json:"openAiServices,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection    `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState              `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess            `json:"publicNetworkAccess,omitempty"`
	StorageServices            *StorageServicesForPatchRequest `json:"storageServices,omitempty"`
	TenantId                   *string                         `json:"tenantId,omitempty"`
}
