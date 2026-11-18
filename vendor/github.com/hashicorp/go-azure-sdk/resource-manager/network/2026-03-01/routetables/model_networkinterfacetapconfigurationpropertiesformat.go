package routetables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceTapConfigurationPropertiesFormat struct {
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	VirtualNetworkTap *VirtualNetworkTap `json:"virtualNetworkTap,omitempty"`
}
