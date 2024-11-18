package backuppolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupPolicyProperties struct {
	BackupPolicyId       *string          `json:"backupPolicyId,omitempty"`
	DailyBackupsToKeep   *int64           `json:"dailyBackupsToKeep,omitempty"`
	Enabled              *bool            `json:"enabled,omitempty"`
	MonthlyBackupsToKeep *int64           `json:"monthlyBackupsToKeep,omitempty"`
	ProvisioningState    *string          `json:"provisioningState,omitempty"`
	VolumeBackups        *[]VolumeBackups `json:"volumeBackups,omitempty"`
	VolumesAssigned      *int64           `json:"volumesAssigned,omitempty"`
	WeeklyBackupsToKeep  *int64           `json:"weeklyBackupsToKeep,omitempty"`
}
