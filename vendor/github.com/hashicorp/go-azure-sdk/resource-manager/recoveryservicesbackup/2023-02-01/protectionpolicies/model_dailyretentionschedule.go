package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DailyRetentionSchedule struct {
	RetentionDuration *RetentionDuration `json:"retentionDuration,omitempty"`
	RetentionTimes    *[]string          `json:"retentionTimes,omitempty"`
}
