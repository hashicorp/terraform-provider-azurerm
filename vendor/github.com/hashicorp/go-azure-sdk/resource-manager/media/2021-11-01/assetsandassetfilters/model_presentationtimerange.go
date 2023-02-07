package assetsandassetfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PresentationTimeRange struct {
	EndTimestamp               *int64 `json:"endTimestamp,omitempty"`
	ForceEndTimestamp          *bool  `json:"forceEndTimestamp,omitempty"`
	LiveBackoffDuration        *int64 `json:"liveBackoffDuration,omitempty"`
	PresentationWindowDuration *int64 `json:"presentationWindowDuration,omitempty"`
	StartTimestamp             *int64 `json:"startTimestamp,omitempty"`
	Timescale                  *int64 `json:"timescale,omitempty"`
}
