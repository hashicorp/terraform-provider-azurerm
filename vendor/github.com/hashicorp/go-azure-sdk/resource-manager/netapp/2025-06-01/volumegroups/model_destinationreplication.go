package volumegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DestinationReplication struct {
	Region          *string          `json:"region,omitempty"`
	ReplicationType *ReplicationType `json:"replicationType,omitempty"`
	ResourceId      *string          `json:"resourceId,omitempty"`
	Zone            *string          `json:"zone,omitempty"`
}
