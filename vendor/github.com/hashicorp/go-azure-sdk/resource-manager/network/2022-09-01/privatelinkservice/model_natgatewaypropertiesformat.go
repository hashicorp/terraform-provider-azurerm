package privatelinkservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NatGatewayPropertiesFormat struct {
	IdleTimeoutInMinutes *int64             `json:"idleTimeoutInMinutes,omitempty"`
	ProvisioningState    *ProvisioningState `json:"provisioningState,omitempty"`
	PublicIPAddresses    *[]SubResource     `json:"publicIpAddresses,omitempty"`
	PublicIPPrefixes     *[]SubResource     `json:"publicIpPrefixes,omitempty"`
	ResourceGuid         *string            `json:"resourceGuid,omitempty"`
	Subnets              *[]SubResource     `json:"subnets,omitempty"`
}
