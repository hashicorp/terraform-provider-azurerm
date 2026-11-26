package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetPublicIPAddressConfigurationProperties struct {
	DeleteOption           *DeleteOptions                                                 `json:"deleteOption,omitempty"`
	DnsSettings            *VirtualMachineScaleSetPublicIPAddressConfigurationDnsSettings `json:"dnsSettings,omitempty"`
	IPTags                 *[]VirtualMachineScaleSetIPTag                                 `json:"ipTags,omitempty"`
	IdleTimeoutInMinutes   *int64                                                         `json:"idleTimeoutInMinutes,omitempty"`
	PublicIPAddressVersion *IPVersion                                                     `json:"publicIPAddressVersion,omitempty"`
	PublicIPPrefix         *SubResource                                                   `json:"publicIPPrefix,omitempty"`
}
