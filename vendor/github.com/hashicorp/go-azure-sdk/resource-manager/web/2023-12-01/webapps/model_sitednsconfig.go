package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteDnsConfig struct {
	DnsAltServer           *string   `json:"dnsAltServer,omitempty"`
	DnsLegacySortOrder     *bool     `json:"dnsLegacySortOrder,omitempty"`
	DnsMaxCacheTimeout     *int64    `json:"dnsMaxCacheTimeout,omitempty"`
	DnsRetryAttemptCount   *int64    `json:"dnsRetryAttemptCount,omitempty"`
	DnsRetryAttemptTimeout *int64    `json:"dnsRetryAttemptTimeout,omitempty"`
	DnsServers             *[]string `json:"dnsServers,omitempty"`
}
