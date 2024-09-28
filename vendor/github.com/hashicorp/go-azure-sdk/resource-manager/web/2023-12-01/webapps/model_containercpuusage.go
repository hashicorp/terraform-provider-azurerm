package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerCPUUsage struct {
	KernelModeUsage *int64   `json:"kernelModeUsage,omitempty"`
	PerCPUUsage     *[]int64 `json:"perCpuUsage,omitempty"`
	TotalUsage      *int64   `json:"totalUsage,omitempty"`
	UserModeUsage   *int64   `json:"userModeUsage,omitempty"`
}
