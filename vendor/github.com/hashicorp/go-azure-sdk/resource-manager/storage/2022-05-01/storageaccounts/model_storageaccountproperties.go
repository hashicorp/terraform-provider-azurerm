package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountProperties struct {
	AccessTier                            *AccessTier                            `json:"accessTier,omitempty"`
	AllowBlobPublicAccess                 *bool                                  `json:"allowBlobPublicAccess,omitempty"`
	AllowCrossTenantReplication           *bool                                  `json:"allowCrossTenantReplication,omitempty"`
	AllowSharedKeyAccess                  *bool                                  `json:"allowSharedKeyAccess,omitempty"`
	AllowedCopyScope                      *AllowedCopyScope                      `json:"allowedCopyScope,omitempty"`
	AzureFilesIdentityBasedAuthentication *AzureFilesIdentityBasedAuthentication `json:"azureFilesIdentityBasedAuthentication"`
	BlobRestoreStatus                     *BlobRestoreStatus                     `json:"blobRestoreStatus"`
	CreationTime                          *string                                `json:"creationTime,omitempty"`
	CustomDomain                          *CustomDomain                          `json:"customDomain"`
	DefaultToOAuthAuthentication          *bool                                  `json:"defaultToOAuthAuthentication,omitempty"`
	DnsEndpointType                       *DnsEndpointType                       `json:"dnsEndpointType,omitempty"`
	Encryption                            *Encryption                            `json:"encryption"`
	FailoverInProgress                    *bool                                  `json:"failoverInProgress,omitempty"`
	GeoReplicationStats                   *GeoReplicationStats                   `json:"geoReplicationStats"`
	ImmutableStorageWithVersioning        *ImmutableStorageAccount               `json:"immutableStorageWithVersioning"`
	IsHnsEnabled                          *bool                                  `json:"isHnsEnabled,omitempty"`
	IsLocalUserEnabled                    *bool                                  `json:"isLocalUserEnabled,omitempty"`
	IsNfsV3Enabled                        *bool                                  `json:"isNfsV3Enabled,omitempty"`
	IsSftpEnabled                         *bool                                  `json:"isSftpEnabled,omitempty"`
	KeyCreationTime                       *KeyCreationTime                       `json:"keyCreationTime"`
	KeyPolicy                             *KeyPolicy                             `json:"keyPolicy"`
	LargeFileSharesState                  *LargeFileSharesState                  `json:"largeFileSharesState,omitempty"`
	LastGeoFailoverTime                   *string                                `json:"lastGeoFailoverTime,omitempty"`
	MinimumTlsVersion                     *MinimumTlsVersion                     `json:"minimumTlsVersion,omitempty"`
	NetworkAcls                           *NetworkRuleSet                        `json:"networkAcls"`
	PrimaryEndpoints                      *Endpoints                             `json:"primaryEndpoints"`
	PrimaryLocation                       *string                                `json:"primaryLocation,omitempty"`
	PrivateEndpointConnections            *[]PrivateEndpointConnection           `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                     *ProvisioningState                     `json:"provisioningState,omitempty"`
	PublicNetworkAccess                   *PublicNetworkAccess                   `json:"publicNetworkAccess,omitempty"`
	RoutingPreference                     *RoutingPreference                     `json:"routingPreference"`
	SasPolicy                             *SasPolicy                             `json:"sasPolicy"`
	SecondaryEndpoints                    *Endpoints                             `json:"secondaryEndpoints"`
	SecondaryLocation                     *string                                `json:"secondaryLocation,omitempty"`
	StatusOfPrimary                       *AccountStatus                         `json:"statusOfPrimary,omitempty"`
	StatusOfSecondary                     *AccountStatus                         `json:"statusOfSecondary,omitempty"`
	StorageAccountSkuConversionStatus     *StorageAccountSkuConversionStatus     `json:"storageAccountSkuConversionStatus"`
	SupportsHTTPSTrafficOnly              *bool                                  `json:"supportsHttpsTrafficOnly,omitempty"`
}

func (o *StorageAccountProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StorageAccountProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *StorageAccountProperties) GetLastGeoFailoverTimeAsTime() (*time.Time, error) {
	if o.LastGeoFailoverTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastGeoFailoverTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StorageAccountProperties) SetLastGeoFailoverTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastGeoFailoverTime = &formatted
}
