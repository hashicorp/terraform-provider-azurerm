package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetspacePropertiesThroughputPoolConfiguration struct {
	DedicatedRUs     *int64 `json:"dedicatedRUs,omitempty"`
	MaxConsumableRUs *int64 `json:"maxConsumableRUs,omitempty"`
	MaxThroughput    *int64 `json:"maxThroughput,omitempty"`
	MinThroughput    *int64 `json:"minThroughput,omitempty"`
}
