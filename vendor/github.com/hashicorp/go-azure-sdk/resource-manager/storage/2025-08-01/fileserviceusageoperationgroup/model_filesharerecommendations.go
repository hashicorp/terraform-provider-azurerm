package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileShareRecommendations struct {
	BandwidthScalar        *float64 `json:"bandwidthScalar,omitempty"`
	BaseBandwidthMiBPerSec *int64   `json:"baseBandwidthMiBPerSec,omitempty"`
	BaseIOPS               *int64   `json:"baseIOPS,omitempty"`
	IoScalar               *float64 `json:"ioScalar,omitempty"`
}
