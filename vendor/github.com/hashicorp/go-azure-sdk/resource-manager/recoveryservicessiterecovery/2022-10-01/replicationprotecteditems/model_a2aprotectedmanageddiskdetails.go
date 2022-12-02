package replicationprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type A2AProtectedManagedDiskDetails struct {
	AllowedDiskLevelOperation              *[]string `json:"allowedDiskLevelOperation,omitempty"`
	DataPendingAtSourceAgentInMB           *float64  `json:"dataPendingAtSourceAgentInMB,omitempty"`
	DataPendingInStagingStorageAccountInMB *float64  `json:"dataPendingInStagingStorageAccountInMB,omitempty"`
	DekKeyVaultArmId                       *string   `json:"dekKeyVaultArmId,omitempty"`
	DiskCapacityInBytes                    *int64    `json:"diskCapacityInBytes,omitempty"`
	DiskId                                 *string   `json:"diskId,omitempty"`
	DiskName                               *string   `json:"diskName,omitempty"`
	DiskState                              *string   `json:"diskState,omitempty"`
	DiskType                               *string   `json:"diskType,omitempty"`
	FailoverDiskName                       *string   `json:"failoverDiskName,omitempty"`
	IsDiskEncrypted                        *bool     `json:"isDiskEncrypted,omitempty"`
	IsDiskKeyEncrypted                     *bool     `json:"isDiskKeyEncrypted,omitempty"`
	KekKeyVaultArmId                       *string   `json:"kekKeyVaultArmId,omitempty"`
	KeyIdentifier                          *string   `json:"keyIdentifier,omitempty"`
	MonitoringJobType                      *string   `json:"monitoringJobType,omitempty"`
	MonitoringPercentageCompletion         *int64    `json:"monitoringPercentageCompletion,omitempty"`
	PrimaryDiskEncryptionSetId             *string   `json:"primaryDiskEncryptionSetId,omitempty"`
	PrimaryStagingAzureStorageAccountId    *string   `json:"primaryStagingAzureStorageAccountId,omitempty"`
	RecoveryDiskEncryptionSetId            *string   `json:"recoveryDiskEncryptionSetId,omitempty"`
	RecoveryOrignalTargetDiskId            *string   `json:"recoveryOrignalTargetDiskId,omitempty"`
	RecoveryReplicaDiskAccountType         *string   `json:"recoveryReplicaDiskAccountType,omitempty"`
	RecoveryReplicaDiskId                  *string   `json:"recoveryReplicaDiskId,omitempty"`
	RecoveryResourceGroupId                *string   `json:"recoveryResourceGroupId,omitempty"`
	RecoveryTargetDiskAccountType          *string   `json:"recoveryTargetDiskAccountType,omitempty"`
	RecoveryTargetDiskId                   *string   `json:"recoveryTargetDiskId,omitempty"`
	ResyncRequired                         *bool     `json:"resyncRequired,omitempty"`
	SecretIdentifier                       *string   `json:"secretIdentifier,omitempty"`
	TfoDiskName                            *string   `json:"tfoDiskName,omitempty"`
}
