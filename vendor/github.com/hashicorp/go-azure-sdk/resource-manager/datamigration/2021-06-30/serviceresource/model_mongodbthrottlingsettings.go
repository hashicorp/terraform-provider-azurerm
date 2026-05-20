package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbThrottlingSettings struct {
	MaxParallelism  *int64 `json:"maxParallelism,omitempty"`
	MinFreeCPU      *int64 `json:"minFreeCpu,omitempty"`
	MinFreeMemoryMb *int64 `json:"minFreeMemoryMb,omitempty"`
}
