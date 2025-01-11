package backuppolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeBackups struct {
	BackupsCount     *int64  `json:"backupsCount,omitempty"`
	PolicyEnabled    *bool   `json:"policyEnabled,omitempty"`
	VolumeName       *string `json:"volumeName,omitempty"`
	VolumeResourceId *string `json:"volumeResourceId,omitempty"`
}
