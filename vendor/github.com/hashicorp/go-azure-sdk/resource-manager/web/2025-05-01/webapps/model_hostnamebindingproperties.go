package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostNameBindingProperties struct {
	AzureResourceName           *string                      `json:"azureResourceName,omitempty"`
	AzureResourceType           *AzureResourceType           `json:"azureResourceType,omitempty"`
	CustomHostNameDnsRecordType *CustomHostNameDnsRecordType `json:"customHostNameDnsRecordType,omitempty"`
	DomainId                    *string                      `json:"domainId,omitempty"`
	HostNameType                *HostNameType                `json:"hostNameType,omitempty"`
	SiteName                    *string                      `json:"siteName,omitempty"`
	SslState                    *SslState                    `json:"sslState,omitempty"`
	Thumbprint                  *string                      `json:"thumbprint,omitempty"`
	VirtualIP                   *string                      `json:"virtualIP,omitempty"`
}
