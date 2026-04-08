package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceRequests struct {
	Cpu        float64      `json:"cpu"`
	Gpu        *GpuResource `json:"gpu,omitempty"`
	MemoryInGB float64      `json:"memoryInGB"`
}
