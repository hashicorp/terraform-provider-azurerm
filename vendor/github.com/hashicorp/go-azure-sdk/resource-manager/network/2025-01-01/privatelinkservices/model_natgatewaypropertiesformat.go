package privatelinkservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatGatewayPropertiesFormat struct {
	IdleTimeoutInMinutes *int64             `json:"idleTimeoutInMinutes,omitempty"`
	ProvisioningState    *ProvisioningState `json:"provisioningState,omitempty"`
	PublicIPAddresses    *[]SubResource     `json:"publicIpAddresses,omitempty"`
	PublicIPAddressesV6  *[]SubResource     `json:"publicIpAddressesV6,omitempty"`
	PublicIPPrefixes     *[]SubResource     `json:"publicIpPrefixes,omitempty"`
	PublicIPPrefixesV6   *[]SubResource     `json:"publicIpPrefixesV6,omitempty"`
	ResourceGuid         *string            `json:"resourceGuid,omitempty"`
	SourceVirtualNetwork *SubResource       `json:"sourceVirtualNetwork,omitempty"`
	Subnets              *[]SubResource     `json:"subnets,omitempty"`
}
