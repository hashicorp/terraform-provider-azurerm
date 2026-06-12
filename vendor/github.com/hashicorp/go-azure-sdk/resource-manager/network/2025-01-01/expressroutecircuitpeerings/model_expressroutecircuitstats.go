package expressroutecircuitpeerings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitStats struct {
	PrimarybytesIn    *int64 `json:"primarybytesIn,omitempty"`
	PrimarybytesOut   *int64 `json:"primarybytesOut,omitempty"`
	SecondarybytesIn  *int64 `json:"secondarybytesIn,omitempty"`
	SecondarybytesOut *int64 `json:"secondarybytesOut,omitempty"`
}
