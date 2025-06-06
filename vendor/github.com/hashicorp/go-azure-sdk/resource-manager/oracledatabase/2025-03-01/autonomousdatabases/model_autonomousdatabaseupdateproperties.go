package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseUpdateProperties struct {
	AdminPassword                        *string                            `json:"adminPassword,omitempty"`
	AutonomousMaintenanceScheduleType    *AutonomousMaintenanceScheduleType `json:"autonomousMaintenanceScheduleType,omitempty"`
	BackupRetentionPeriodInDays          *int64                             `json:"backupRetentionPeriodInDays,omitempty"`
	ComputeCount                         *float64                           `json:"computeCount,omitempty"`
	CpuCoreCount                         *int64                             `json:"cpuCoreCount,omitempty"`
	CustomerContacts                     *[]CustomerContact                 `json:"customerContacts,omitempty"`
	DataStorageSizeInGbs                 *int64                             `json:"dataStorageSizeInGbs,omitempty"`
	DataStorageSizeInTbs                 *int64                             `json:"dataStorageSizeInTbs,omitempty"`
	DatabaseEdition                      *DatabaseEditionType               `json:"databaseEdition,omitempty"`
	DisplayName                          *string                            `json:"displayName,omitempty"`
	IsAutoScalingEnabled                 *bool                              `json:"isAutoScalingEnabled,omitempty"`
	IsAutoScalingForStorageEnabled       *bool                              `json:"isAutoScalingForStorageEnabled,omitempty"`
	IsLocalDataGuardEnabled              *bool                              `json:"isLocalDataGuardEnabled,omitempty"`
	IsMtlsConnectionRequired             *bool                              `json:"isMtlsConnectionRequired,omitempty"`
	LicenseModel                         *LicenseModel                      `json:"licenseModel,omitempty"`
	LocalAdgAutoFailoverMaxDataLossLimit *int64                             `json:"localAdgAutoFailoverMaxDataLossLimit,omitempty"`
	LongTermBackupSchedule               *LongTermBackUpScheduleDetails     `json:"longTermBackupSchedule,omitempty"`
	OpenMode                             *OpenModeType                      `json:"openMode,omitempty"`
	PeerDbId                             *string                            `json:"peerDbId,omitempty"`
	PermissionLevel                      *PermissionLevelType               `json:"permissionLevel,omitempty"`
	Role                                 *RoleType                          `json:"role,omitempty"`
	ScheduledOperations                  *ScheduledOperationsTypeUpdate     `json:"scheduledOperations,omitempty"`
	WhitelistedIPs                       *[]string                          `json:"whitelistedIps,omitempty"`
}
