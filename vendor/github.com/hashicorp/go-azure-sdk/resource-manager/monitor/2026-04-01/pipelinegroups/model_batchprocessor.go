package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchProcessor struct {
	BatchSize *int64 `json:"batchSize,omitempty"`
	Timeout   *int64 `json:"timeout,omitempty"`
}
