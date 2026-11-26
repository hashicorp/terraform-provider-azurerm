package privateendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendAddressPoolPropertiesFormat struct {
	BackendAddresses        *[]ApplicationGatewayBackendAddress `json:"backendAddresses,omitempty"`
	BackendIPConfigurations *[]NetworkInterfaceIPConfiguration  `json:"backendIPConfigurations,omitempty"`
	ProvisioningState       *ProvisioningState                  `json:"provisioningState,omitempty"`
}
