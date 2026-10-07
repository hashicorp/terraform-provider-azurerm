package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficDetectionRule struct {
	PacketsPerSecond *int64           `json:"packetsPerSecond,omitempty"`
	TrafficType      *DdosTrafficType `json:"trafficType,omitempty"`
}
