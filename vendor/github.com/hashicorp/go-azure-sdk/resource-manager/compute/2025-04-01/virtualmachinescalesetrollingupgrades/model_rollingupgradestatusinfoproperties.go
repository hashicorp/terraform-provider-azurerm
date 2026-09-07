package virtualmachinescalesetrollingupgrades

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollingUpgradeStatusInfoProperties struct {
	Error         *ApiError                    `json:"error,omitempty"`
	Policy        *RollingUpgradePolicy        `json:"policy,omitempty"`
	Progress      *RollingUpgradeProgressInfo  `json:"progress,omitempty"`
	RunningStatus *RollingUpgradeRunningStatus `json:"runningStatus,omitempty"`
}
