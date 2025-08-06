package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AProtectedDiskDetails struct {
	AllowedDiskLevelOperation              *[]string `json:"allowedDiskLevelOperation,omitempty"`
	DataPendingAtSourceAgentInMB           *float64  `json:"dataPendingAtSourceAgentInMB,omitempty"`
	DataPendingInStagingStorageAccountInMB *float64  `json:"dataPendingInStagingStorageAccountInMB,omitempty"`
	DekKeyVaultArmId                       *string   `json:"dekKeyVaultArmId,omitempty"`
	DiskCapacityInBytes                    *int64    `json:"diskCapacityInBytes,omitempty"`
	DiskName                               *string   `json:"diskName,omitempty"`
	DiskState                              *string   `json:"diskState,omitempty"`
	DiskType                               *string   `json:"diskType,omitempty"`
	DiskUri                                *string   `json:"diskUri,omitempty"`
	FailoverDiskName                       *string   `json:"failoverDiskName,omitempty"`
	IsDiskEncrypted                        *bool     `json:"isDiskEncrypted,omitempty"`
	IsDiskKeyEncrypted                     *bool     `json:"isDiskKeyEncrypted,omitempty"`
	KekKeyVaultArmId                       *string   `json:"kekKeyVaultArmId,omitempty"`
	KeyIdentifier                          *string   `json:"keyIdentifier,omitempty"`
	MonitoringJobType                      *string   `json:"monitoringJobType,omitempty"`
	MonitoringPercentageCompletion         *int64    `json:"monitoringPercentageCompletion,omitempty"`
	PrimaryDiskAzureStorageAccountId       *string   `json:"primaryDiskAzureStorageAccountId,omitempty"`
	PrimaryStagingAzureStorageAccountId    *string   `json:"primaryStagingAzureStorageAccountId,omitempty"`
	RecoveryAzureStorageAccountId          *string   `json:"recoveryAzureStorageAccountId,omitempty"`
	RecoveryDiskUri                        *string   `json:"recoveryDiskUri,omitempty"`
	ResyncRequired                         *bool     `json:"resyncRequired,omitempty"`
	SecretIdentifier                       *string   `json:"secretIdentifier,omitempty"`
	TfoDiskName                            *string   `json:"tfoDiskName,omitempty"`
}
