package storageaccounts

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoReplicationStats struct {
	CanFailover                   *bool                          `json:"canFailover,omitempty"`
	CanPlannedFailover            *bool                          `json:"canPlannedFailover,omitempty"`
	LastSyncTime                  *string                        `json:"lastSyncTime,omitempty"`
	PostFailoverRedundancy        *PostFailoverRedundancy        `json:"postFailoverRedundancy,omitempty"`
	PostPlannedFailoverRedundancy *PostPlannedFailoverRedundancy `json:"postPlannedFailoverRedundancy,omitempty"`
	Status                        *GeoReplicationStatus          `json:"status,omitempty"`
}

func (o *GeoReplicationStats) GetLastSyncTimeAsTime() (*time.Time, error) {
	if o.LastSyncTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GeoReplicationStats) SetLastSyncTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTime = &formatted
}
