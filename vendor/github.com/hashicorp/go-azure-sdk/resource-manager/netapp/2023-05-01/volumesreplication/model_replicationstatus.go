package volumesreplication

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationStatus struct {
	ErrorMessage       *string             `json:"errorMessage,omitempty"`
	Healthy            *bool               `json:"healthy,omitempty"`
	MirrorState        *MirrorState        `json:"mirrorState,omitempty"`
	RelationshipStatus *RelationshipStatus `json:"relationshipStatus,omitempty"`
	TotalProgress      *string             `json:"totalProgress,omitempty"`
}
