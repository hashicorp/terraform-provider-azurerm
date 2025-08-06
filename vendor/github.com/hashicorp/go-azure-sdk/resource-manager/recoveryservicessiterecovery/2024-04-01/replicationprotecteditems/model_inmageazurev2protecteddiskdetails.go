package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageAzureV2ProtectedDiskDetails struct {
	DiskCapacityInBytes                 *int64   `json:"diskCapacityInBytes,omitempty"`
	DiskId                              *string  `json:"diskId,omitempty"`
	DiskName                            *string  `json:"diskName,omitempty"`
	DiskResized                         *string  `json:"diskResized,omitempty"`
	FileSystemCapacityInBytes           *int64   `json:"fileSystemCapacityInBytes,omitempty"`
	HealthErrorCode                     *string  `json:"healthErrorCode,omitempty"`
	LastRpoCalculatedTime               *string  `json:"lastRpoCalculatedTime,omitempty"`
	ProgressHealth                      *string  `json:"progressHealth,omitempty"`
	ProgressStatus                      *string  `json:"progressStatus,omitempty"`
	ProtectionStage                     *string  `json:"protectionStage,omitempty"`
	PsDataInMegaBytes                   *float64 `json:"psDataInMegaBytes,omitempty"`
	ResyncDurationInSeconds             *int64   `json:"resyncDurationInSeconds,omitempty"`
	ResyncLast15MinutesTransferredBytes *int64   `json:"resyncLast15MinutesTransferredBytes,omitempty"`
	ResyncLastDataTransferTimeUTC       *string  `json:"resyncLastDataTransferTimeUTC,omitempty"`
	ResyncProcessedBytes                *int64   `json:"resyncProcessedBytes,omitempty"`
	ResyncProgressPercentage            *int64   `json:"resyncProgressPercentage,omitempty"`
	ResyncRequired                      *string  `json:"resyncRequired,omitempty"`
	ResyncStartTime                     *string  `json:"resyncStartTime,omitempty"`
	ResyncTotalTransferredBytes         *int64   `json:"resyncTotalTransferredBytes,omitempty"`
	RpoInSeconds                        *int64   `json:"rpoInSeconds,omitempty"`
	SecondsToTakeSwitchProvider         *int64   `json:"secondsToTakeSwitchProvider,omitempty"`
	SourceDataInMegaBytes               *float64 `json:"sourceDataInMegaBytes,omitempty"`
	TargetDataInMegaBytes               *float64 `json:"targetDataInMegaBytes,omitempty"`
}

func (o *InMageAzureV2ProtectedDiskDetails) GetLastRpoCalculatedTimeAsTime() (*time.Time, error) {
	if o.LastRpoCalculatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRpoCalculatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageAzureV2ProtectedDiskDetails) SetLastRpoCalculatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRpoCalculatedTime = &formatted
}

func (o *InMageAzureV2ProtectedDiskDetails) GetResyncLastDataTransferTimeUTCAsTime() (*time.Time, error) {
	if o.ResyncLastDataTransferTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ResyncLastDataTransferTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageAzureV2ProtectedDiskDetails) SetResyncLastDataTransferTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ResyncLastDataTransferTimeUTC = &formatted
}

func (o *InMageAzureV2ProtectedDiskDetails) GetResyncStartTimeAsTime() (*time.Time, error) {
	if o.ResyncStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ResyncStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageAzureV2ProtectedDiskDetails) SetResyncStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ResyncStartTime = &formatted
}
