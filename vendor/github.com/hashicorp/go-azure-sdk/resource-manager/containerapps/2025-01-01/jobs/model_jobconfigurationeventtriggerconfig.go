package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobConfigurationEventTriggerConfig struct {
	Parallelism            *int64    `json:"parallelism,omitempty"`
	ReplicaCompletionCount *int64    `json:"replicaCompletionCount,omitempty"`
	Scale                  *JobScale `json:"scale,omitempty"`
}
