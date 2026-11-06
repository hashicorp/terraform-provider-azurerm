package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlatformProfile struct {
	ContainerRegistry       *ContainerRegistryProfile      `json:"containerRegistry,omitempty"`
	IssuerURL               string                         `json:"issuerUrl"`
	ManagedResourceGroup    *string                        `json:"managedResourceGroup,omitempty"`
	NetworkSecurityGroupId  string                         `json:"networkSecurityGroupId"`
	OperatorsAuthentication OperatorsAuthenticationProfile `json:"operatorsAuthentication"`
	OutboundType            *OutboundType                  `json:"outboundType,omitempty"`
	SubnetId                string                         `json:"subnetId"`
	VnetIntegrationSubnetId string                         `json:"vnetIntegrationSubnetId"`
}
