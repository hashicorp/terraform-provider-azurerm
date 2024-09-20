package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetInfo struct {
	CertBlob       *string      `json:"certBlob,omitempty"`
	CertThumbprint *string      `json:"certThumbprint,omitempty"`
	DnsServers     *string      `json:"dnsServers,omitempty"`
	IsSwift        *bool        `json:"isSwift,omitempty"`
	ResyncRequired *bool        `json:"resyncRequired,omitempty"`
	Routes         *[]VnetRoute `json:"routes,omitempty"`
	VnetResourceId *string      `json:"vnetResourceId,omitempty"`
}
