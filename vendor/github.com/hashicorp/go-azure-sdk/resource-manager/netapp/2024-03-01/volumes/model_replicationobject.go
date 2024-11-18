package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationObject struct {
	EndpointType           *EndpointType        `json:"endpointType,omitempty"`
	RemoteVolumeRegion     *string              `json:"remoteVolumeRegion,omitempty"`
	RemoteVolumeResourceId string               `json:"remoteVolumeResourceId"`
	ReplicationId          *string              `json:"replicationId,omitempty"`
	ReplicationSchedule    *ReplicationSchedule `json:"replicationSchedule,omitempty"`
}
