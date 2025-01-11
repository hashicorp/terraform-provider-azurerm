package publicipprefixes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPPrefixPropertiesFormat struct {
	CustomIPPrefix                      *SubResource                 `json:"customIPPrefix,omitempty"`
	IPPrefix                            *string                      `json:"ipPrefix,omitempty"`
	IPTags                              *[]IPTag                     `json:"ipTags,omitempty"`
	LoadBalancerFrontendIPConfiguration *SubResource                 `json:"loadBalancerFrontendIpConfiguration,omitempty"`
	NatGateway                          *NatGateway                  `json:"natGateway,omitempty"`
	PrefixLength                        *int64                       `json:"prefixLength,omitempty"`
	ProvisioningState                   *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicIPAddressVersion              *IPVersion                   `json:"publicIPAddressVersion,omitempty"`
	PublicIPAddresses                   *[]ReferencedPublicIPAddress `json:"publicIPAddresses,omitempty"`
	ResourceGuid                        *string                      `json:"resourceGuid,omitempty"`
}
