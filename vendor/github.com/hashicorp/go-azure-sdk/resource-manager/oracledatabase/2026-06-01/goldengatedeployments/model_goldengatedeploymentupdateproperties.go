package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GoldenGateDeploymentUpdateProperties struct {
	BackupSchedule           *BackupScheduleType           `json:"backupSchedule,omitempty"`
	CpuCoreCount             *int64                        `json:"cpuCoreCount,omitempty"`
	LicenseModel             *LicenseModel                 `json:"licenseModel,omitempty"`
	MaintenanceConfiguration *MaintenanceConfigurationType `json:"maintenanceConfiguration,omitempty"`
	MaintenanceWindow        *MaintenanceWindowType        `json:"maintenanceWindow,omitempty"`
}
