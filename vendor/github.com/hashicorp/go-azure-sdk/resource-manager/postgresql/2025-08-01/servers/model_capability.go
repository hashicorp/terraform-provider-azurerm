package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Capability struct {
	FastProvisioningSupported            *FastProvisioningSupport                                              `json:"fastProvisioningSupported,omitempty"`
	GeoBackupSupported                   *GeographicallyRedundantBackupSupport                                 `json:"geoBackupSupported,omitempty"`
	Name                                 *string                                                               `json:"name,omitempty"`
	OnlineResizeSupported                *OnlineStorageResizeSupport                                           `json:"onlineResizeSupported,omitempty"`
	Reason                               *string                                                               `json:"reason,omitempty"`
	Restricted                           *LocationRestricted                                                   `json:"restricted,omitempty"`
	Status                               *CapabilityStatus                                                     `json:"status,omitempty"`
	StorageAutoGrowthSupported           *StorageAutoGrowthSupport                                             `json:"storageAutoGrowthSupported,omitempty"`
	SupportedFastProvisioningEditions    *[]FastProvisioningEditionCapability                                  `json:"supportedFastProvisioningEditions,omitempty"`
	SupportedFeatures                    *[]SupportedFeature                                                   `json:"supportedFeatures,omitempty"`
	SupportedServerEditions              *[]ServerEditionCapability                                            `json:"supportedServerEditions,omitempty"`
	SupportedServerVersions              *[]ServerVersionCapability                                            `json:"supportedServerVersions,omitempty"`
	ZoneRedundantHaAndGeoBackupSupported *ZoneRedundantHighAvailabilityAndGeographicallyRedundantBackupSupport `json:"zoneRedundantHaAndGeoBackupSupported,omitempty"`
	ZoneRedundantHaSupported             *ZoneRedundantHighAvailabilitySupport                                 `json:"zoneRedundantHaSupported,omitempty"`
}
