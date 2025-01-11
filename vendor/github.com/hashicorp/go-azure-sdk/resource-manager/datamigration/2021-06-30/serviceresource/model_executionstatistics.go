package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutionStatistics struct {
	CpuTimeMs      *float64                   `json:"cpuTimeMs,omitempty"`
	ElapsedTimeMs  *float64                   `json:"elapsedTimeMs,omitempty"`
	ExecutionCount *int64                     `json:"executionCount,omitempty"`
	HasErrors      *bool                      `json:"hasErrors,omitempty"`
	SqlErrors      *[]string                  `json:"sqlErrors,omitempty"`
	WaitStats      *map[string]WaitStatistics `json:"waitStats,omitempty"`
}
