package integrationaccountbatchconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchReleaseCriteria struct {
	BatchSize    *int64                     `json:"batchSize,omitempty"`
	MessageCount *int64                     `json:"messageCount,omitempty"`
	Recurrence   *WorkflowTriggerRecurrence `json:"recurrence,omitempty"`
}
