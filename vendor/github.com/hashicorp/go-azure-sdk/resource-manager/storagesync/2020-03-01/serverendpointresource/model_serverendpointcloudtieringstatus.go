package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointCloudTieringStatus struct {
	CachePerformance            *CloudTieringCachePerformance            `json:"cachePerformance,omitempty"`
	DatePolicyStatus            *CloudTieringDatePolicyStatus            `json:"datePolicyStatus,omitempty"`
	FilesNotTiering             *CloudTieringFilesNotTiering             `json:"filesNotTiering,omitempty"`
	Health                      *ServerEndpointCloudTieringHealthState   `json:"health,omitempty"`
	HealthLastUpdatedTimestamp  *string                                  `json:"healthLastUpdatedTimestamp,omitempty"`
	LastCloudTieringResult      *int64                                   `json:"lastCloudTieringResult,omitempty"`
	LastSuccessTimestamp        *string                                  `json:"lastSuccessTimestamp,omitempty"`
	LastUpdatedTimestamp        *string                                  `json:"lastUpdatedTimestamp,omitempty"`
	SpaceSavings                *CloudTieringSpaceSavings                `json:"spaceSavings,omitempty"`
	VolumeFreeSpacePolicyStatus *CloudTieringVolumeFreeSpacePolicyStatus `json:"volumeFreeSpacePolicyStatus,omitempty"`
}

func (o *ServerEndpointCloudTieringStatus) GetHealthLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.HealthLastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.HealthLastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointCloudTieringStatus) SetHealthLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.HealthLastUpdatedTimestamp = &formatted
}

func (o *ServerEndpointCloudTieringStatus) GetLastSuccessTimestampAsTime() (*time.Time, error) {
	if o.LastSuccessTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointCloudTieringStatus) SetLastSuccessTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessTimestamp = &formatted
}

func (o *ServerEndpointCloudTieringStatus) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointCloudTieringStatus) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
