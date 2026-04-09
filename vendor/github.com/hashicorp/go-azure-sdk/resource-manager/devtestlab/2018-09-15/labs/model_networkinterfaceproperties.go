package labs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceProperties struct {
	DnsName                            *string                             `json:"dnsName,omitempty"`
	PrivateIPAddress                   *string                             `json:"privateIpAddress,omitempty"`
	PublicIPAddress                    *string                             `json:"publicIpAddress,omitempty"`
	PublicIPAddressId                  *string                             `json:"publicIpAddressId,omitempty"`
	RdpAuthority                       *string                             `json:"rdpAuthority,omitempty"`
	SharedPublicIPAddressConfiguration *SharedPublicIPAddressConfiguration `json:"sharedPublicIpAddressConfiguration,omitempty"`
	SshAuthority                       *string                             `json:"sshAuthority,omitempty"`
	SubnetId                           *string                             `json:"subnetId,omitempty"`
	VirtualNetworkId                   *string                             `json:"virtualNetworkId,omitempty"`
}
