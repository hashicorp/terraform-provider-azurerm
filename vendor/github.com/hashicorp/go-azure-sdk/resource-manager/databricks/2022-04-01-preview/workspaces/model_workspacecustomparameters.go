package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCustomParameters struct {
	AmlWorkspaceId                  *WorkspaceCustomStringParameter  `json:"amlWorkspaceId"`
	CustomPrivateSubnetName         *WorkspaceCustomStringParameter  `json:"customPrivateSubnetName"`
	CustomPublicSubnetName          *WorkspaceCustomStringParameter  `json:"customPublicSubnetName"`
	CustomVirtualNetworkId          *WorkspaceCustomStringParameter  `json:"customVirtualNetworkId"`
	EnableNoPublicIP                *WorkspaceCustomBooleanParameter `json:"enableNoPublicIp"`
	Encryption                      *WorkspaceEncryptionParameter    `json:"encryption"`
	LoadBalancerBackendPoolName     *WorkspaceCustomStringParameter  `json:"loadBalancerBackendPoolName"`
	LoadBalancerId                  *WorkspaceCustomStringParameter  `json:"loadBalancerId"`
	NatGatewayName                  *WorkspaceCustomStringParameter  `json:"natGatewayName"`
	PrepareEncryption               *WorkspaceCustomBooleanParameter `json:"prepareEncryption"`
	PublicIPName                    *WorkspaceCustomStringParameter  `json:"publicIpName"`
	RequireInfrastructureEncryption *WorkspaceCustomBooleanParameter `json:"requireInfrastructureEncryption"`
	ResourceTags                    *WorkspaceCustomObjectParameter  `json:"resourceTags"`
	StorageAccountName              *WorkspaceCustomStringParameter  `json:"storageAccountName"`
	StorageAccountSkuName           *WorkspaceCustomStringParameter  `json:"storageAccountSkuName"`
	VnetAddressPrefix               *WorkspaceCustomStringParameter  `json:"vnetAddressPrefix"`
}
