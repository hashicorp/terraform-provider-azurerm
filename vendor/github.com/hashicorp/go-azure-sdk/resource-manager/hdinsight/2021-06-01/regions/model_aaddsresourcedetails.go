package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AaddsResourceDetails struct {
	DomainName                     *string `json:"domainName,omitempty"`
	InitialSyncComplete            *bool   `json:"initialSyncComplete,omitempty"`
	LdapsEnabled                   *bool   `json:"ldapsEnabled,omitempty"`
	LdapsPublicCertificateInBase64 *string `json:"ldapsPublicCertificateInBase64,omitempty"`
	ResourceId                     *string `json:"resourceId,omitempty"`
	SubnetId                       *string `json:"subnetId,omitempty"`
	TenantId                       *string `json:"tenantId,omitempty"`
}
