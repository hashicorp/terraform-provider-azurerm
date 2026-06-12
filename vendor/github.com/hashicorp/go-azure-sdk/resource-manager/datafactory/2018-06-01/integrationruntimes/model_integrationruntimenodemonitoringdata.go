package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeNodeMonitoringData struct {
	AvailableMemoryInMB   *int64   `json:"availableMemoryInMB,omitempty"`
	ConcurrentJobsLimit   *int64   `json:"concurrentJobsLimit,omitempty"`
	ConcurrentJobsRunning *int64   `json:"concurrentJobsRunning,omitempty"`
	CpuUtilization        *int64   `json:"cpuUtilization,omitempty"`
	MaxConcurrentJobs     *int64   `json:"maxConcurrentJobs,omitempty"`
	NodeName              *string  `json:"nodeName,omitempty"`
	ReceivedBytes         *float64 `json:"receivedBytes,omitempty"`
	SentBytes             *float64 `json:"sentBytes,omitempty"`
}
