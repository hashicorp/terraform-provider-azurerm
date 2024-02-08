package serverbasedperformancetier

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerformanceTierServiceLevelObjectives struct {
	Edition                *string `json:"edition,omitempty"`
	HardwareGeneration     *string `json:"hardwareGeneration,omitempty"`
	Id                     *string `json:"id,omitempty"`
	MaxBackupRetentionDays *int64  `json:"maxBackupRetentionDays,omitempty"`
	MaxStorageMB           *int64  `json:"maxStorageMB,omitempty"`
	MinBackupRetentionDays *int64  `json:"minBackupRetentionDays,omitempty"`
	MinStorageMB           *int64  `json:"minStorageMB,omitempty"`
	VCore                  *int64  `json:"vCore,omitempty"`
}
