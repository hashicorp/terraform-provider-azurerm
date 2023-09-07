package serverbasedperformancetier

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerformanceTierProperties struct {
	Id                     *string                                  `json:"id,omitempty"`
	MaxBackupRetentionDays *int64                                   `json:"maxBackupRetentionDays,omitempty"`
	MaxLargeStorageMB      *int64                                   `json:"maxLargeStorageMB,omitempty"`
	MaxStorageMB           *int64                                   `json:"maxStorageMB,omitempty"`
	MinBackupRetentionDays *int64                                   `json:"minBackupRetentionDays,omitempty"`
	MinLargeStorageMB      *int64                                   `json:"minLargeStorageMB,omitempty"`
	MinStorageMB           *int64                                   `json:"minStorageMB,omitempty"`
	ServiceLevelObjectives *[]PerformanceTierServiceLevelObjectives `json:"serviceLevelObjectives,omitempty"`
}
