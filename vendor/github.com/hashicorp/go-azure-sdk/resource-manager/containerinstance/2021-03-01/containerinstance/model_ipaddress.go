package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IpAddress struct {
	DnsNameLabel *string                     `json:"dnsNameLabel,omitempty"`
	Fqdn         *string                     `json:"fqdn,omitempty"`
	Ip           *string                     `json:"ip,omitempty"`
	Ports        []Port                      `json:"ports"`
	Type         ContainerGroupIpAddressType `json:"type"`
}
