package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Traces struct {
	CacheSizeMegabytes  *int64           `json:"cacheSizeMegabytes,omitempty"`
	Mode                *OperationalMode `json:"mode,omitempty"`
	SelfTracing         *SelfTracing     `json:"selfTracing,omitempty"`
	SpanChannelCapacity *int64           `json:"spanChannelCapacity,omitempty"`
}
