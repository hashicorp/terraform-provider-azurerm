package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SiteLimits struct {
	MaxDiskSizeInMb  *int64   `json:"maxDiskSizeInMb,omitempty"`
	MaxMemoryInMb    *int64   `json:"maxMemoryInMb,omitempty"`
	MaxPercentageCPU *float64 `json:"maxPercentageCpu,omitempty"`
}
