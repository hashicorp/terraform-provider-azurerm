package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WindowsGmsaProfile struct {
	DnsServer      *string `json:"dnsServer,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty"`
	RootDomainName *string `json:"rootDomainName,omitempty"`
}
