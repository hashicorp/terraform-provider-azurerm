package routetables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfaceDnsSettings struct {
	AppliedDnsServers        *[]string `json:"appliedDnsServers,omitempty"`
	DnsServers               *[]string `json:"dnsServers,omitempty"`
	InternalDnsNameLabel     *string   `json:"internalDnsNameLabel,omitempty"`
	InternalDomainNameSuffix *string   `json:"internalDomainNameSuffix,omitempty"`
	InternalFqdn             *string   `json:"internalFqdn,omitempty"`
}
