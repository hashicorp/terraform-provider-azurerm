package backups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupStatus struct {
	ErrorMessage          *string             `json:"errorMessage,omitempty"`
	Healthy               *bool               `json:"healthy,omitempty"`
	LastTransferSize      *int64              `json:"lastTransferSize,omitempty"`
	LastTransferType      *string             `json:"lastTransferType,omitempty"`
	MirrorState           *MirrorState        `json:"mirrorState,omitempty"`
	RelationshipStatus    *RelationshipStatus `json:"relationshipStatus,omitempty"`
	TotalTransferBytes    *int64              `json:"totalTransferBytes,omitempty"`
	TransferProgressBytes *int64              `json:"transferProgressBytes,omitempty"`
	UnhealthyReason       *string             `json:"unhealthyReason,omitempty"`
}
