package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerCPUStatistics struct {
	CpuUsage       *ContainerCPUUsage       `json:"cpuUsage,omitempty"`
	OnlineCPUCount *int64                   `json:"onlineCpuCount,omitempty"`
	SystemCPUUsage *int64                   `json:"systemCpuUsage,omitempty"`
	ThrottlingData *ContainerThrottlingData `json:"throttlingData,omitempty"`
}
