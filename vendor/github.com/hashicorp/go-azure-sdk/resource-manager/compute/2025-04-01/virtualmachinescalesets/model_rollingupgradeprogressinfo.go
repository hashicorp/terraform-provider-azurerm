package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollingUpgradeProgressInfo struct {
	FailedInstanceCount     *int64 `json:"failedInstanceCount,omitempty"`
	InProgressInstanceCount *int64 `json:"inProgressInstanceCount,omitempty"`
	PendingInstanceCount    *int64 `json:"pendingInstanceCount,omitempty"`
	SuccessfulInstanceCount *int64 `json:"successfulInstanceCount,omitempty"`
}
