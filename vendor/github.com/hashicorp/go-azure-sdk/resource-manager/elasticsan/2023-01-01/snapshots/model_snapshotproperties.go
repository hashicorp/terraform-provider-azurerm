package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotProperties struct {
	CreationData        SnapshotCreationData `json:"creationData"`
	ProvisioningState   *ProvisioningStates  `json:"provisioningState,omitempty"`
	SourceVolumeSizeGiB *int64               `json:"sourceVolumeSizeGiB,omitempty"`
	VolumeName          *string              `json:"volumeName,omitempty"`
}
