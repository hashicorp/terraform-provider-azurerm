package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OsPatchingStatus struct {
	LatestPatchTime     *string          `json:"latestPatchTime,omitempty"`
	OsPatchingErrors    *[]ErrorResponse `json:"osPatchingErrors,omitempty"`
	PatchStatus         *PatchStatus     `json:"patchStatus,omitempty"`
	RebootPending       *bool            `json:"rebootPending,omitempty"`
	ScheduledRebootTime *string          `json:"scheduledRebootTime,omitempty"`
}
