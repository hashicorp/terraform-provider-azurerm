package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapabilitiesCapabilities struct {
	MaxCPU        *float64 `json:"maxCpu,omitempty"`
	MaxGpuCount   *float64 `json:"maxGpuCount,omitempty"`
	MaxMemoryInGB *float64 `json:"maxMemoryInGB,omitempty"`
}
