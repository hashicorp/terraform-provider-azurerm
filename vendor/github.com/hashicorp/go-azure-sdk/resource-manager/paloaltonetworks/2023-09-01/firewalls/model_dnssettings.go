package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DNSSettings struct {
	DnsServers     *[]IPAddress    `json:"dnsServers,omitempty"`
	EnableDnsProxy *DNSProxy       `json:"enableDnsProxy,omitempty"`
	EnabledDnsType *EnabledDNSType `json:"enabledDnsType,omitempty"`
}
