package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateProfileRollingUpdateProfile struct {
	InPlaceUpdate           *bool   `json:"inPlaceUpdate,omitempty"`
	MaxBatchPercent         *int64  `json:"maxBatchPercent,omitempty"`
	MaxUnhealthyPercent     *int64  `json:"maxUnhealthyPercent,omitempty"`
	PauseTimeBetweenBatches *string `json:"pauseTimeBetweenBatches,omitempty"`
}
