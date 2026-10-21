package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsSummary struct {
	FailedJobs     *int64 `json:"failedJobs,omitempty"`
	InProgressJobs *int64 `json:"inProgressJobs,omitempty"`
	SuspendedJobs  *int64 `json:"suspendedJobs,omitempty"`
}
