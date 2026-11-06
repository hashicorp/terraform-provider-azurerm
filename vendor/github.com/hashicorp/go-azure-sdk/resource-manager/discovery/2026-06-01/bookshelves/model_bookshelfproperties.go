package bookshelves

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BookshelfProperties struct {
	BookshelfUri                   *string                          `json:"bookshelfUri,omitempty"`
	CustomerManagedKeys            *CustomerManagedKeys             `json:"customerManagedKeys,omitempty"`
	KeyVaultProperties             *BookshelfKeyVaultProperties     `json:"keyVaultProperties,omitempty"`
	LogAnalyticsClusterId          *string                          `json:"logAnalyticsClusterId,omitempty"`
	ManagedOnBehalfOfConfiguration *WithMoboBrokerResources         `json:"managedOnBehalfOfConfiguration,omitempty"`
	ManagedResourceGroup           *string                          `json:"managedResourceGroup,omitempty"`
	PrivateEndpointConnections     *[]PrivateEndpointConnection     `json:"privateEndpointConnections,omitempty"`
	PrivateEndpointSubnetId        *string                          `json:"privateEndpointSubnetId,omitempty"`
	ProvisioningState              *ProvisioningState               `json:"provisioningState,omitempty"`
	PublicNetworkAccess            *PublicNetworkAccess             `json:"publicNetworkAccess,omitempty"`
	SearchSubnetId                 *string                          `json:"searchSubnetId,omitempty"`
	WorkloadIdentities             *map[string]UserAssignedIdentity `json:"workloadIdentities,omitempty"`
}
