package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotProperties struct {
	CreationData      *CreationData `json:"creationData"`
	EnableFIPS        *bool         `json:"enableFIPS,omitempty"`
	KubernetesVersion *string       `json:"kubernetesVersion,omitempty"`
	NodeImageVersion  *string       `json:"nodeImageVersion,omitempty"`
	OsSku             *OSSKU        `json:"osSku,omitempty"`
	OsType            *OSType       `json:"osType,omitempty"`
	SnapshotType      *SnapshotType `json:"snapshotType,omitempty"`
	VmSize            *string       `json:"vmSize,omitempty"`
}
