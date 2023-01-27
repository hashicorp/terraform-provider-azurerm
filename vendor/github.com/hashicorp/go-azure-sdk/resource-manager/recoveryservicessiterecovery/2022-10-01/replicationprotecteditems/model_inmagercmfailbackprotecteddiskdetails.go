package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmFailbackProtectedDiskDetails struct {
	CapacityInBytes               *int64                        `json:"capacityInBytes,omitempty"`
	DataPendingAtSourceAgentInMB  *float64                      `json:"dataPendingAtSourceAgentInMB,omitempty"`
	DataPendingInLogDataStoreInMB *float64                      `json:"dataPendingInLogDataStoreInMB,omitempty"`
	DiskId                        *string                       `json:"diskId,omitempty"`
	DiskName                      *string                       `json:"diskName,omitempty"`
	DiskUuid                      *string                       `json:"diskUuid,omitempty"`
	IrDetails                     *InMageRcmFailbackSyncDetails `json:"irDetails,omitempty"`
	IsInitialReplicationComplete  *string                       `json:"isInitialReplicationComplete,omitempty"`
	IsOSDisk                      *string                       `json:"isOSDisk,omitempty"`
	LastSyncTime                  *string                       `json:"lastSyncTime,omitempty"`
	ResyncDetails                 *InMageRcmFailbackSyncDetails `json:"resyncDetails,omitempty"`
}

func (o *InMageRcmFailbackProtectedDiskDetails) GetLastSyncTimeAsTime() (*time.Time, error) {
	if o.LastSyncTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTime, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmFailbackProtectedDiskDetails) SetLastSyncTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTime = &formatted
}
