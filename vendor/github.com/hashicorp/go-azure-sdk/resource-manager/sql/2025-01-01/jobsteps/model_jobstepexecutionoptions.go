package jobsteps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStepExecutionOptions struct {
	InitialRetryIntervalSeconds    *int64   `json:"initialRetryIntervalSeconds,omitempty"`
	MaximumRetryIntervalSeconds    *int64   `json:"maximumRetryIntervalSeconds,omitempty"`
	RetryAttempts                  *int64   `json:"retryAttempts,omitempty"`
	RetryIntervalBackoffMultiplier *float64 `json:"retryIntervalBackoffMultiplier,omitempty"`
	TimeoutSeconds                 *int64   `json:"timeoutSeconds,omitempty"`
}
