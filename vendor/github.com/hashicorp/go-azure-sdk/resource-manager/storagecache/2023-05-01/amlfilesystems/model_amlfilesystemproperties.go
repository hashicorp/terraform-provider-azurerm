package amlfilesystems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemProperties struct {
	ClientInfo                *AmlFilesystemClientInfo                 `json:"clientInfo,omitempty"`
	EncryptionSettings        *AmlFilesystemEncryptionSettings         `json:"encryptionSettings,omitempty"`
	FilesystemSubnet          string                                   `json:"filesystemSubnet"`
	Health                    *AmlFilesystemHealth                     `json:"health,omitempty"`
	Hsm                       *AmlFilesystemPropertiesHsm              `json:"hsm,omitempty"`
	MaintenanceWindow         AmlFilesystemPropertiesMaintenanceWindow `json:"maintenanceWindow"`
	ProvisioningState         *AmlFilesystemProvisioningStateType      `json:"provisioningState,omitempty"`
	StorageCapacityTiB        float64                                  `json:"storageCapacityTiB"`
	ThroughputProvisionedMBps *int64                                   `json:"throughputProvisionedMBps,omitempty"`
}
