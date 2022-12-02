package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageDiskExclusionInput struct {
	DiskSignatureOptions *[]InMageDiskSignatureExclusionOptions `json:"diskSignatureOptions,omitempty"`
	VolumeOptions        *[]InMageVolumeExclusionOptions        `json:"volumeOptions,omitempty"`
}
