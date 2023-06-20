package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionVolume struct {
	CapacityInBytes     *int64  `json:"capacityInBytes,omitempty"`
	FreeSpaceInBytes    *int64  `json:"freeSpaceInBytes,omitempty"`
	ThresholdPercentage *int64  `json:"thresholdPercentage,omitempty"`
	VolumeName          *string `json:"volumeName,omitempty"`
}
