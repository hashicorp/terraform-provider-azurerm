package cloudexadatainfrastructures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudExadataInfrastructureProperties struct {
	ActivatedStorageCount       *int64                                    `json:"activatedStorageCount,omitempty"`
	AdditionalStorageCount      *int64                                    `json:"additionalStorageCount,omitempty"`
	AvailableStorageSizeInGbs   *int64                                    `json:"availableStorageSizeInGbs,omitempty"`
	ComputeCount                *int64                                    `json:"computeCount,omitempty"`
	CpuCount                    *int64                                    `json:"cpuCount,omitempty"`
	CustomerContacts            *[]CustomerContact                        `json:"customerContacts,omitempty"`
	DataStorageSizeInTbs        *float64                                  `json:"dataStorageSizeInTbs,omitempty"`
	DbNodeStorageSizeInGbs      *int64                                    `json:"dbNodeStorageSizeInGbs,omitempty"`
	DbServerVersion             *string                                   `json:"dbServerVersion,omitempty"`
	DisplayName                 string                                    `json:"displayName"`
	EstimatedPatchingTime       *EstimatedPatchingTime                    `json:"estimatedPatchingTime,omitempty"`
	LastMaintenanceRunId        *string                                   `json:"lastMaintenanceRunId,omitempty"`
	LifecycleDetails            *string                                   `json:"lifecycleDetails,omitempty"`
	LifecycleState              *CloudExadataInfrastructureLifecycleState `json:"lifecycleState,omitempty"`
	MaintenanceWindow           *MaintenanceWindow                        `json:"maintenanceWindow,omitempty"`
	MaxCPUCount                 *int64                                    `json:"maxCpuCount,omitempty"`
	MaxDataStorageInTbs         *float64                                  `json:"maxDataStorageInTbs,omitempty"`
	MaxDbNodeStorageSizeInGbs   *int64                                    `json:"maxDbNodeStorageSizeInGbs,omitempty"`
	MaxMemoryInGbs              *int64                                    `json:"maxMemoryInGbs,omitempty"`
	MemorySizeInGbs             *int64                                    `json:"memorySizeInGbs,omitempty"`
	MonthlyDbServerVersion      *string                                   `json:"monthlyDbServerVersion,omitempty"`
	MonthlyStorageServerVersion *string                                   `json:"monthlyStorageServerVersion,omitempty"`
	NextMaintenanceRunId        *string                                   `json:"nextMaintenanceRunId,omitempty"`
	OciURL                      *string                                   `json:"ociUrl,omitempty"`
	Ocid                        *string                                   `json:"ocid,omitempty"`
	ProvisioningState           *AzureResourceProvisioningState           `json:"provisioningState,omitempty"`
	Shape                       string                                    `json:"shape"`
	StorageCount                *int64                                    `json:"storageCount,omitempty"`
	StorageServerVersion        *string                                   `json:"storageServerVersion,omitempty"`
	TimeCreated                 *string                                   `json:"timeCreated,omitempty"`
	TotalStorageSizeInGbs       *int64                                    `json:"totalStorageSizeInGbs,omitempty"`
}
