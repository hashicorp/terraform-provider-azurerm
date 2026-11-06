package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentProperties struct {
	BackupSchedule            *BackupScheduleType             `json:"backupSchedule,omitempty"`
	Category                  *CategoryType                   `json:"category,omitempty"`
	Compartment               *string                         `json:"compartment,omitempty"`
	CpuCoreCount              *int64                          `json:"cpuCoreCount,omitempty"`
	DeploymentType            *DeploymentType                 `json:"deploymentType,omitempty"`
	DeploymentURL             *string                         `json:"deploymentUrl,omitempty"`
	DisplayName               string                          `json:"displayName"`
	EnvironmentType           *SetupType                      `json:"environmentType,omitempty"`
	IngressIPs                *[]string                       `json:"ingressIps,omitempty"`
	IsAutoScalingEnabled      *bool                           `json:"isAutoScalingEnabled,omitempty"`
	IsPublic                  *bool                           `json:"isPublic,omitempty"`
	LicenseModel              *LicenseModel                   `json:"licenseModel,omitempty"`
	LifecycleDetails          *string                         `json:"lifecycleDetails,omitempty"`
	LifecycleState            *DeploymentLifecycleState       `json:"lifecycleState,omitempty"`
	MaintenanceConfiguration  *MaintenanceConfigurationType   `json:"maintenanceConfiguration,omitempty"`
	MaintenanceWindow         *MaintenanceWindowType          `json:"maintenanceWindow,omitempty"`
	NetworkAnchorId           string                          `json:"networkAnchorId"`
	Ocid                      *string                         `json:"ocid,omitempty"`
	OggData                   *OggDeploymentDetails           `json:"oggData,omitempty"`
	PrivateIPAddress          *string                         `json:"privateIpAddress,omitempty"`
	ProvisioningState         *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
	ResourceAnchorId          string                          `json:"resourceAnchorId"`
	StorageUtilizationInBytes *int64                          `json:"storageUtilizationInBytes,omitempty"`
	TimeCreated               *string                         `json:"timeCreated,omitempty"`
	TimeUpdated               *string                         `json:"timeUpdated,omitempty"`
	TimeZone                  *string                         `json:"timeZone,omitempty"`
	Version                   *string                         `json:"version,omitempty"`
}
