package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationObject struct {
	DestinationReplications        *[]DestinationReplication            `json:"destinationReplications,omitempty"`
	EndpointType                   *EndpointType                        `json:"endpointType,omitempty"`
	ExternalReplicationSetupInfo   *string                              `json:"externalReplicationSetupInfo,omitempty"`
	ExternalReplicationSetupStatus *ExternalReplicationSetupStatus      `json:"externalReplicationSetupStatus,omitempty"`
	MirrorState                    *MirrorState                         `json:"mirrorState,omitempty"`
	RelationshipStatus             *VolumeReplicationRelationshipStatus `json:"relationshipStatus,omitempty"`
	RemotePath                     *RemotePath                          `json:"remotePath,omitempty"`
	RemoteVolumeRegion             *string                              `json:"remoteVolumeRegion,omitempty"`
	RemoteVolumeResourceId         *string                              `json:"remoteVolumeResourceId,omitempty"`
	ReplicationId                  *string                              `json:"replicationId,omitempty"`
	ReplicationSchedule            *ReplicationSchedule                 `json:"replicationSchedule,omitempty"`
}
