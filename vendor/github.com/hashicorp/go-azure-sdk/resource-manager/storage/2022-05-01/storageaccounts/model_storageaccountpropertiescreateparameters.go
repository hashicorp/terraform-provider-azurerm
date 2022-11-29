package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountPropertiesCreateParameters struct {
	AccessTier                            *AccessTier                            `json:"accessTier,omitempty"`
	AllowBlobPublicAccess                 *bool                                  `json:"allowBlobPublicAccess,omitempty"`
	AllowCrossTenantReplication           *bool                                  `json:"allowCrossTenantReplication,omitempty"`
	AllowSharedKeyAccess                  *bool                                  `json:"allowSharedKeyAccess,omitempty"`
	AllowedCopyScope                      *AllowedCopyScope                      `json:"allowedCopyScope,omitempty"`
	AzureFilesIdentityBasedAuthentication *AzureFilesIdentityBasedAuthentication `json:"azureFilesIdentityBasedAuthentication"`
	CustomDomain                          *CustomDomain                          `json:"customDomain"`
	DefaultToOAuthAuthentication          *bool                                  `json:"defaultToOAuthAuthentication,omitempty"`
	DnsEndpointType                       *DnsEndpointType                       `json:"dnsEndpointType,omitempty"`
	Encryption                            *Encryption                            `json:"encryption"`
	ImmutableStorageWithVersioning        *ImmutableStorageAccount               `json:"immutableStorageWithVersioning"`
	IsHnsEnabled                          *bool                                  `json:"isHnsEnabled,omitempty"`
	IsLocalUserEnabled                    *bool                                  `json:"isLocalUserEnabled,omitempty"`
	IsNfsV3Enabled                        *bool                                  `json:"isNfsV3Enabled,omitempty"`
	IsSftpEnabled                         *bool                                  `json:"isSftpEnabled,omitempty"`
	KeyPolicy                             *KeyPolicy                             `json:"keyPolicy"`
	LargeFileSharesState                  *LargeFileSharesState                  `json:"largeFileSharesState,omitempty"`
	MinimumTlsVersion                     *MinimumTlsVersion                     `json:"minimumTlsVersion,omitempty"`
	NetworkAcls                           *NetworkRuleSet                        `json:"networkAcls"`
	PublicNetworkAccess                   *PublicNetworkAccess                   `json:"publicNetworkAccess,omitempty"`
	RoutingPreference                     *RoutingPreference                     `json:"routingPreference"`
	SasPolicy                             *SasPolicy                             `json:"sasPolicy"`
	SupportsHTTPSTrafficOnly              *bool                                  `json:"supportsHttpsTrafficOnly,omitempty"`
}
