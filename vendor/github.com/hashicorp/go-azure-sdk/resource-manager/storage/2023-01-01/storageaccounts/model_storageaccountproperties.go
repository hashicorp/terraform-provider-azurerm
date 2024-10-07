package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountProperties struct {
	AccessTier                            *AccessTier                            `json:"accessTier,omitempty"`
	AccountMigrationInProgress            *bool                                  `json:"accountMigrationInProgress,omitempty"`
	AllowBlobPublicAccess                 *bool                                  `json:"allowBlobPublicAccess,omitempty"`
	AllowCrossTenantReplication           *bool                                  `json:"allowCrossTenantReplication,omitempty"`
	AllowSharedKeyAccess                  *bool                                  `json:"allowSharedKeyAccess,omitempty"`
	AllowedCopyScope                      *AllowedCopyScope                      `json:"allowedCopyScope,omitempty"`
	AzureFilesIdentityBasedAuthentication *AzureFilesIdentityBasedAuthentication `json:"azureFilesIdentityBasedAuthentication,omitempty"`
	BlobRestoreStatus                     *BlobRestoreStatus                     `json:"blobRestoreStatus,omitempty"`
	CreationTime                          *string                                `json:"creationTime,omitempty"`
	CustomDomain                          *CustomDomain                          `json:"customDomain,omitempty"`
	DefaultToOAuthAuthentication          *bool                                  `json:"defaultToOAuthAuthentication,omitempty"`
	DnsEndpointType                       *DnsEndpointType                       `json:"dnsEndpointType,omitempty"`
	Encryption                            *Encryption                            `json:"encryption,omitempty"`
	FailoverInProgress                    *bool                                  `json:"failoverInProgress,omitempty"`
	GeoReplicationStats                   *GeoReplicationStats                   `json:"geoReplicationStats,omitempty"`
	ImmutableStorageWithVersioning        *ImmutableStorageAccount               `json:"immutableStorageWithVersioning,omitempty"`
	IsHnsEnabled                          *bool                                  `json:"isHnsEnabled,omitempty"`
	IsLocalUserEnabled                    *bool                                  `json:"isLocalUserEnabled,omitempty"`
	IsNfsV3Enabled                        *bool                                  `json:"isNfsV3Enabled,omitempty"`
	IsSftpEnabled                         *bool                                  `json:"isSftpEnabled,omitempty"`
	IsSkuConversionBlocked                *bool                                  `json:"isSkuConversionBlocked,omitempty"`
	KeyCreationTime                       *KeyCreationTime                       `json:"keyCreationTime,omitempty"`
	KeyPolicy                             *KeyPolicy                             `json:"keyPolicy,omitempty"`
	LargeFileSharesState                  *LargeFileSharesState                  `json:"largeFileSharesState,omitempty"`
	LastGeoFailoverTime                   *string                                `json:"lastGeoFailoverTime,omitempty"`
	MinimumTlsVersion                     *MinimumTlsVersion                     `json:"minimumTlsVersion,omitempty"`
	NetworkAcls                           *NetworkRuleSet                        `json:"networkAcls,omitempty"`
	PrimaryEndpoints                      *Endpoints                             `json:"primaryEndpoints,omitempty"`
	PrimaryLocation                       *string                                `json:"primaryLocation,omitempty"`
	PrivateEndpointConnections            *[]PrivateEndpointConnection           `json:"privateEndpointConnections,omitempty"`
	ProvisioningState                     *ProvisioningState                     `json:"provisioningState,omitempty"`
	PublicNetworkAccess                   *PublicNetworkAccess                   `json:"publicNetworkAccess,omitempty"`
	RoutingPreference                     *RoutingPreference                     `json:"routingPreference,omitempty"`
	SasPolicy                             *SasPolicy                             `json:"sasPolicy,omitempty"`
	SecondaryEndpoints                    *Endpoints                             `json:"secondaryEndpoints,omitempty"`
	SecondaryLocation                     *string                                `json:"secondaryLocation,omitempty"`
	StatusOfPrimary                       *AccountStatus                         `json:"statusOfPrimary,omitempty"`
	StatusOfSecondary                     *AccountStatus                         `json:"statusOfSecondary,omitempty"`
	StorageAccountSkuConversionStatus     *StorageAccountSkuConversionStatus     `json:"storageAccountSkuConversionStatus,omitempty"`
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
