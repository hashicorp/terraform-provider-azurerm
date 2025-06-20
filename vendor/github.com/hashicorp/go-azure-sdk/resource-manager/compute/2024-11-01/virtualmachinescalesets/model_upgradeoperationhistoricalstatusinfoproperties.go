package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradeOperationHistoricalStatusInfoProperties struct {
	Error                *ApiError                      `json:"error,omitempty"`
	Progress             *RollingUpgradeProgressInfo    `json:"progress,omitempty"`
	RollbackInfo         *RollbackStatusInfo            `json:"rollbackInfo,omitempty"`
	RunningStatus        *UpgradeOperationHistoryStatus `json:"runningStatus,omitempty"`
	StartedBy            *UpgradeOperationInvoker       `json:"startedBy,omitempty"`
	TargetImageReference *ImageReference                `json:"targetImageReference,omitempty"`
}
