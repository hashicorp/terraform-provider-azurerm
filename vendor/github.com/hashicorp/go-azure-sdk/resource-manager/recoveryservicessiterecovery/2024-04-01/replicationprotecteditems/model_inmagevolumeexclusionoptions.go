package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageVolumeExclusionOptions struct {
	OnlyExcludeIfSingleVolume *string `json:"onlyExcludeIfSingleVolume,omitempty"`
	VolumeLabel               *string `json:"volumeLabel,omitempty"`
}
