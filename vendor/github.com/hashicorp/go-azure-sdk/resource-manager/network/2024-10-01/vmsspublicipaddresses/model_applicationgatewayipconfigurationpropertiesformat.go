package vmsspublicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayIPConfigurationPropertiesFormat struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Subnet            *SubResource       `json:"subnet,omitempty"`
}
