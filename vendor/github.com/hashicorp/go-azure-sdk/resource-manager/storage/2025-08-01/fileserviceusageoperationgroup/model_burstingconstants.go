package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BurstingConstants struct {
	BurstFloorIOPS        *int64   `json:"burstFloorIOPS,omitempty"`
	BurstIOScalar         *float64 `json:"burstIOScalar,omitempty"`
	BurstTimeframeSeconds *int64   `json:"burstTimeframeSeconds,omitempty"`
}
