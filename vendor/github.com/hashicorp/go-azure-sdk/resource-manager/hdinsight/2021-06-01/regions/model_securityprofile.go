package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityProfile struct {
	AaddsResourceId      *string        `json:"aaddsResourceId,omitempty"`
	ClusterUsersGroupDNs *[]string      `json:"clusterUsersGroupDNs,omitempty"`
	DirectoryType        *DirectoryType `json:"directoryType,omitempty"`
	Domain               *string        `json:"domain,omitempty"`
	DomainUserPassword   *string        `json:"domainUserPassword,omitempty"`
	DomainUsername       *string        `json:"domainUsername,omitempty"`
	LdapsURLs            *[]string      `json:"ldapsUrls,omitempty"`
	MsiResourceId        *string        `json:"msiResourceId,omitempty"`
	OrganizationalUnitDN *string        `json:"organizationalUnitDN,omitempty"`
}
