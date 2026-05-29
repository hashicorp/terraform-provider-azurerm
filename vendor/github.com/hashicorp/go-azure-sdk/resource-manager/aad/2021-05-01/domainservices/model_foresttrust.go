package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForestTrust struct {
	FriendlyName      *string `json:"friendlyName,omitempty"`
	RemoteDnsIPs      *string `json:"remoteDnsIps,omitempty"`
	TrustDirection    *string `json:"trustDirection,omitempty"`
	TrustPassword     *string `json:"trustPassword,omitempty"`
	TrustedDomainFqdn *string `json:"trustedDomainFqdn,omitempty"`
}
