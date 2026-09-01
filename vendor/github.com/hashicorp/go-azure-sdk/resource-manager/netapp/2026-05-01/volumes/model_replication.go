package volumes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Replication struct {
	EndpointType            *EndpointType           `json:"endpointType,omitempty"`
	MirrorState             *ReplicationMirrorState `json:"mirrorState,omitempty"`
	RemoteVolumeRegion      *string                 `json:"remoteVolumeRegion,omitempty"`
	RemoteVolumeResourceId  *string                 `json:"remoteVolumeResourceId,omitempty"`
	ReplicationCreationTime *string                 `json:"replicationCreationTime,omitempty"`
	ReplicationDeletionTime *string                 `json:"replicationDeletionTime,omitempty"`
	ReplicationId           *string                 `json:"replicationId,omitempty"`
	ReplicationSchedule     *ReplicationSchedule    `json:"replicationSchedule,omitempty"`
}

func (o *Replication) GetReplicationCreationTimeAsTime() (*time.Time, error) {
	if o.ReplicationCreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ReplicationCreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Replication) SetReplicationCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReplicationCreationTime = &formatted
}

func (o *Replication) GetReplicationDeletionTimeAsTime() (*time.Time, error) {
	if o.ReplicationDeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ReplicationDeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Replication) SetReplicationDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReplicationDeletionTime = &formatted
}
