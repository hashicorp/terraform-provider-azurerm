package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceProperties struct {
	DnsSettings       *InterfaceDNSSettings   `json:"dnsSettings,omitempty"`
	IPConfigurations  *[]IPConfiguration      `json:"ipConfigurations,omitempty"`
	MacAddress        *string                 `json:"macAddress,omitempty"`
	ProvisioningState *ProvisioningStateEnum  `json:"provisioningState,omitempty"`
	Status            *NetworkInterfaceStatus `json:"status,omitempty"`
}
