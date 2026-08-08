package backupresourcevaultconfigs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupResourceVaultConfig struct {
	EnhancedSecurityState            *EnhancedSecurityState  `json:"enhancedSecurityState,omitempty"`
	IsSoftDeleteFeatureStateEditable *bool                   `json:"isSoftDeleteFeatureStateEditable,omitempty"`
	ResourceGuardOperationRequests   *[]string               `json:"resourceGuardOperationRequests,omitempty"`
	SoftDeleteFeatureState           *SoftDeleteFeatureState `json:"softDeleteFeatureState,omitempty"`
	SoftDeleteRetentionPeriodInDays  *int64                  `json:"softDeleteRetentionPeriodInDays,omitempty"`
	StorageModelType                 *StorageType            `json:"storageModelType,omitempty"`
	StorageType                      *StorageType            `json:"storageType,omitempty"`
	StorageTypeState                 *StorageTypeState       `json:"storageTypeState,omitempty"`
}
