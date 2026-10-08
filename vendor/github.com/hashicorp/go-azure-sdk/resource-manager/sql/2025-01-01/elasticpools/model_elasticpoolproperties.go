package elasticpools

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticPoolProperties struct {
	AutoPauseDelay               *int64                          `json:"autoPauseDelay,omitempty"`
	AvailabilityZone             *AvailabilityZoneType           `json:"availabilityZone,omitempty"`
	CreationDate                 *string                         `json:"creationDate,omitempty"`
	HighAvailabilityReplicaCount *int64                          `json:"highAvailabilityReplicaCount,omitempty"`
	LicenseType                  *ElasticPoolLicenseType         `json:"licenseType,omitempty"`
	MaintenanceConfigurationId   *string                         `json:"maintenanceConfigurationId,omitempty"`
	MaxSizeBytes                 *int64                          `json:"maxSizeBytes,omitempty"`
	MinCapacity                  *float64                        `json:"minCapacity,omitempty"`
	PerDatabaseSettings          *ElasticPoolPerDatabaseSettings `json:"perDatabaseSettings,omitempty"`
	PreferredEnclaveType         *AlwaysEncryptedEnclaveType     `json:"preferredEnclaveType,omitempty"`
	State                        *ElasticPoolState               `json:"state,omitempty"`
	ZoneRedundant                *bool                           `json:"zoneRedundant,omitempty"`
}

func (o *ElasticPoolProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ElasticPoolProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
