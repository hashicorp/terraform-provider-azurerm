package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheActiveDirectorySettings struct {
	CacheNetBiosName      string                                   `json:"cacheNetBiosName"`
	Credentials           *CacheActiveDirectorySettingsCredentials `json:"credentials,omitempty"`
	DomainJoined          *DomainJoinedType                        `json:"domainJoined,omitempty"`
	DomainName            string                                   `json:"domainName"`
	DomainNetBiosName     string                                   `json:"domainNetBiosName"`
	PrimaryDnsIPAddress   string                                   `json:"primaryDnsIpAddress"`
	SecondaryDnsIPAddress *string                                  `json:"secondaryDnsIpAddress,omitempty"`
}
