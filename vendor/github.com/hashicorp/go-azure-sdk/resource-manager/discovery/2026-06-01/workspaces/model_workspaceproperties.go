package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProperties struct {
	AgentSubnetId                  *string                      `json:"agentSubnetId,omitempty"`
	CustomerManagedKeys            *CustomerManagedKeys         `json:"customerManagedKeys,omitempty"`
	KeyVaultProperties             *KeyVaultProperties          `json:"keyVaultProperties,omitempty"`
	LogAnalyticsClusterId          *string                      `json:"logAnalyticsClusterId,omitempty"`
	ManagedOnBehalfOfConfiguration *WithMoboBrokerResources     `json:"managedOnBehalfOfConfiguration,omitempty"`
	ManagedResourceGroup           *string                      `json:"managedResourceGroup,omitempty"`
	PrivateEndpointConnections     *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	PrivateEndpointSubnetId        *string                      `json:"privateEndpointSubnetId,omitempty"`
	ProvisioningState              *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess            *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	SupercomputerIds               *[]string                    `json:"supercomputerIds,omitempty"`
	WorkspaceApiUri                *string                      `json:"workspaceApiUri,omitempty"`
	WorkspaceIdentity              Identity                     `json:"workspaceIdentity"`
	WorkspaceSubnetId              *string                      `json:"workspaceSubnetId,omitempty"`
	WorkspaceUiUri                 *string                      `json:"workspaceUiUri,omitempty"`
}
