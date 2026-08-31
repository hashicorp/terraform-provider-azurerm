package elasticpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticPoolUpdateProperties struct {
	AutoPauseDelay               *int64                          `json:"autoPauseDelay,omitempty"`
	AvailabilityZone             *AvailabilityZoneType           `json:"availabilityZone,omitempty"`
	HighAvailabilityReplicaCount *int64                          `json:"highAvailabilityReplicaCount,omitempty"`
	LicenseType                  *ElasticPoolLicenseType         `json:"licenseType,omitempty"`
	MaintenanceConfigurationId   *string                         `json:"maintenanceConfigurationId,omitempty"`
	MaxSizeBytes                 *int64                          `json:"maxSizeBytes,omitempty"`
	MinCapacity                  *float64                        `json:"minCapacity,omitempty"`
	PerDatabaseSettings          *ElasticPoolPerDatabaseSettings `json:"perDatabaseSettings,omitempty"`
	PreferredEnclaveType         *AlwaysEncryptedEnclaveType     `json:"preferredEnclaveType,omitempty"`
	ZoneRedundant                *bool                           `json:"zoneRedundant,omitempty"`
}
