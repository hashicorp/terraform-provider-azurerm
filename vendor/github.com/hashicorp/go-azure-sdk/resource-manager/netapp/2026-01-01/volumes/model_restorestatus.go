package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestoreStatus struct {
	ErrorMessage       *string                          `json:"errorMessage,omitempty"`
	Healthy            *bool                            `json:"healthy,omitempty"`
	MirrorState        *MirrorState                     `json:"mirrorState,omitempty"`
	RelationshipStatus *VolumeRestoreRelationshipStatus `json:"relationshipStatus,omitempty"`
	TotalTransferBytes *int64                           `json:"totalTransferBytes,omitempty"`
	UnhealthyReason    *string                          `json:"unhealthyReason,omitempty"`
}
