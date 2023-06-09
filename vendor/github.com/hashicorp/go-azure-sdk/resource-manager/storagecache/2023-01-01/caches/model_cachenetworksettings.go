package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheNetworkSettings struct {
	DnsSearchDomain  *string   `json:"dnsSearchDomain,omitempty"`
	DnsServers       *[]string `json:"dnsServers,omitempty"`
	Mtu              *int64    `json:"mtu,omitempty"`
	NtpServer        *string   `json:"ntpServer,omitempty"`
	UtilityAddresses *[]string `json:"utilityAddresses,omitempty"`
}
