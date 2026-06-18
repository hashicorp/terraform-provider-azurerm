package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateIntegrationRuntimeNodeRequest struct {
	ConcurrentJobsLimit *int64 `json:"concurrentJobsLimit,omitempty"`
}
