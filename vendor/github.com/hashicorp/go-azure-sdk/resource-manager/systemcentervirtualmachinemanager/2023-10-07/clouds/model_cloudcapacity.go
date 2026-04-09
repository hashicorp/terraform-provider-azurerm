package clouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudCapacity struct {
	CpuCount *int64 `json:"cpuCount,omitempty"`
	MemoryMB *int64 `json:"memoryMB,omitempty"`
	VMCount  *int64 `json:"vmCount,omitempty"`
}
