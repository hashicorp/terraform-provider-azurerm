package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupVault struct {
	FeatureSettings                 *FeatureSettings     `json:"featureSettings,omitempty"`
	IsVaultProtectedByResourceGuard *bool                `json:"isVaultProtectedByResourceGuard,omitempty"`
	MonitoringSettings              *MonitoringSettings  `json:"monitoringSettings,omitempty"`
	ProvisioningState               *ProvisioningState   `json:"provisioningState,omitempty"`
	ResourceMoveDetails             *ResourceMoveDetails `json:"resourceMoveDetails,omitempty"`
	ResourceMoveState               *ResourceMoveState   `json:"resourceMoveState,omitempty"`
	SecureScore                     *SecureScoreLevel    `json:"secureScore,omitempty"`
	SecuritySettings                *SecuritySettings    `json:"securitySettings,omitempty"`
	StorageSettings                 []StorageSetting     `json:"storageSettings"`
}
