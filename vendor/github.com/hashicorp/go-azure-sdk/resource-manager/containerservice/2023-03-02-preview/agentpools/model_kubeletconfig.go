package agentpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KubeletConfig struct {
	AllowedUnsafeSysctls  *[]string `json:"allowedUnsafeSysctls,omitempty"`
	ContainerLogMaxFiles  *int64    `json:"containerLogMaxFiles,omitempty"`
	ContainerLogMaxSizeMB *int64    `json:"containerLogMaxSizeMB,omitempty"`
	CpuCfsQuota           *bool     `json:"cpuCfsQuota,omitempty"`
	CpuCfsQuotaPeriod     *string   `json:"cpuCfsQuotaPeriod,omitempty"`
	CpuManagerPolicy      *string   `json:"cpuManagerPolicy,omitempty"`
	FailSwapOn            *bool     `json:"failSwapOn,omitempty"`
	ImageGcHighThreshold  *int64    `json:"imageGcHighThreshold,omitempty"`
	ImageGcLowThreshold   *int64    `json:"imageGcLowThreshold,omitempty"`
	PodMaxPids            *int64    `json:"podMaxPids,omitempty"`
	TopologyManagerPolicy *string   `json:"topologyManagerPolicy,omitempty"`
}
