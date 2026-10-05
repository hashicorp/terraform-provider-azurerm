package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineNetworkInterfaceConfigurationProperties struct {
	AuxiliaryMode               *NetworkInterfaceAuxiliaryMode                          `json:"auxiliaryMode,omitempty"`
	AuxiliarySku                *NetworkInterfaceAuxiliarySku                           `json:"auxiliarySku,omitempty"`
	DeleteOption                *DeleteOptions                                          `json:"deleteOption,omitempty"`
	DisableTcpStateTracking     *bool                                                   `json:"disableTcpStateTracking,omitempty"`
	DnsSettings                 *VirtualMachineNetworkInterfaceDnsSettingsConfiguration `json:"dnsSettings,omitempty"`
	DscpConfiguration           *SubResource                                            `json:"dscpConfiguration,omitempty"`
	EnableAcceleratedNetworking *bool                                                   `json:"enableAcceleratedNetworking,omitempty"`
	EnableFpga                  *bool                                                   `json:"enableFpga,omitempty"`
	EnableIPForwarding          *bool                                                   `json:"enableIPForwarding,omitempty"`
	IPConfigurations            []VirtualMachineNetworkInterfaceIPConfiguration         `json:"ipConfigurations"`
	NetworkSecurityGroup        *SubResource                                            `json:"networkSecurityGroup,omitempty"`
	Primary                     *bool                                                   `json:"primary,omitempty"`
}
