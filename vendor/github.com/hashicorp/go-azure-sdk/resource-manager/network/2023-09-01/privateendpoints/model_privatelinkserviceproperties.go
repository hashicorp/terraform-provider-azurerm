package privateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceProperties struct {
	Alias                                *string                              `json:"alias,omitempty"`
	AutoApproval                         *ResourceSet                         `json:"autoApproval,omitempty"`
	EnableProxyProtocol                  *bool                                `json:"enableProxyProtocol,omitempty"`
	Fqdns                                *[]string                            `json:"fqdns,omitempty"`
	IPConfigurations                     *[]PrivateLinkServiceIPConfiguration `json:"ipConfigurations,omitempty"`
	LoadBalancerFrontendIPConfigurations *[]FrontendIPConfiguration           `json:"loadBalancerFrontendIpConfigurations,omitempty"`
	NetworkInterfaces                    *[]NetworkInterface                  `json:"networkInterfaces,omitempty"`
	PrivateEndpointConnections           *[]PrivateEndpointConnection         `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                    *ProvisioningState                   `json:"provisioningState,omitempty"`
	Visibility                           *ResourceSet                         `json:"visibility,omitempty"`
}
