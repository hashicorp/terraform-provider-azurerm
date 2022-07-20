package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCustomParameters struct {
	AmlWorkspaceId                  *WorkspaceCustomStringParameter  `json:"amlWorkspaceId,omitempty"`
	CustomPrivateSubnetName         *WorkspaceCustomStringParameter  `json:"customPrivateSubnetName,omitempty"`
	CustomPublicSubnetName          *WorkspaceCustomStringParameter  `json:"customPublicSubnetName,omitempty"`
	CustomVirtualNetworkId          *WorkspaceCustomStringParameter  `json:"customVirtualNetworkId,omitempty"`
	EnableNoPublicIp                *WorkspaceCustomBooleanParameter `json:"enableNoPublicIp,omitempty"`
	Encryption                      *WorkspaceEncryptionParameter    `json:"encryption,omitempty"`
	LoadBalancerBackendPoolName     *WorkspaceCustomStringParameter  `json:"loadBalancerBackendPoolName,omitempty"`
	LoadBalancerId                  *WorkspaceCustomStringParameter  `json:"loadBalancerId,omitempty"`
	NatGatewayName                  *WorkspaceCustomStringParameter  `json:"natGatewayName,omitempty"`
	PrepareEncryption               *WorkspaceCustomBooleanParameter `json:"prepareEncryption,omitempty"`
	PublicIpName                    *WorkspaceCustomStringParameter  `json:"publicIpName,omitempty"`
	RequireInfrastructureEncryption *WorkspaceCustomBooleanParameter `json:"requireInfrastructureEncryption,omitempty"`
	ResourceTags                    *WorkspaceCustomObjectParameter  `json:"resourceTags,omitempty"`
	StorageAccountName              *WorkspaceCustomStringParameter  `json:"storageAccountName,omitempty"`
	StorageAccountSkuName           *WorkspaceCustomStringParameter  `json:"storageAccountSkuName,omitempty"`
	VnetAddressPrefix               *WorkspaceCustomStringParameter  `json:"vnetAddressPrefix,omitempty"`
}
