package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationOperationPredicate struct {
	RemoteVolumeRegion      *string
	RemoteVolumeResourceId  *string
	ReplicationCreationTime *string
	ReplicationDeletionTime *string
	ReplicationId           *string
}

func (p ReplicationOperationPredicate) Matches(input Replication) bool {

	if p.RemoteVolumeRegion != nil && (input.RemoteVolumeRegion == nil || *p.RemoteVolumeRegion != *input.RemoteVolumeRegion) {
		return false
	}

	if p.RemoteVolumeResourceId != nil && (input.RemoteVolumeResourceId == nil || *p.RemoteVolumeResourceId != *input.RemoteVolumeResourceId) {
		return false
	}

	if p.ReplicationCreationTime != nil && (input.ReplicationCreationTime == nil || *p.ReplicationCreationTime != *input.ReplicationCreationTime) {
		return false
	}

	if p.ReplicationDeletionTime != nil && (input.ReplicationDeletionTime == nil || *p.ReplicationDeletionTime != *input.ReplicationDeletionTime) {
		return false
	}

	if p.ReplicationId != nil && (input.ReplicationId == nil || *p.ReplicationId != *input.ReplicationId) {
		return false
	}

	return true
}

type VolumeOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p VolumeOperationPredicate) Matches(input Volume) bool {

	if p.Etag != nil && (input.Etag == nil || *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}
