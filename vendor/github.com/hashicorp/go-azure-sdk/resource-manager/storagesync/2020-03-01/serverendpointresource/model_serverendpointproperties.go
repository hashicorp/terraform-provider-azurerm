package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointProperties struct {
	CloudTiering                                *FeatureStatus                    `json:"cloudTiering,omitempty"`
	CloudTieringStatus                          *ServerEndpointCloudTieringStatus `json:"cloudTieringStatus,omitempty"`
	FriendlyName                                *string                           `json:"friendlyName,omitempty"`
	InitialDownloadPolicy                       *InitialDownloadPolicy            `json:"initialDownloadPolicy,omitempty"`
	LastOperationName                           *string                           `json:"lastOperationName,omitempty"`
	LastWorkflowId                              *string                           `json:"lastWorkflowId,omitempty"`
	LocalCacheMode                              *LocalCacheMode                   `json:"localCacheMode,omitempty"`
	OfflineDataTransfer                         *FeatureStatus                    `json:"offlineDataTransfer,omitempty"`
	OfflineDataTransferShareName                *string                           `json:"offlineDataTransferShareName,omitempty"`
	OfflineDataTransferStorageAccountResourceId *string                           `json:"offlineDataTransferStorageAccountResourceId,omitempty"`
	OfflineDataTransferStorageAccountTenantId   *string                           `json:"offlineDataTransferStorageAccountTenantId,omitempty"`
	ProvisioningState                           *string                           `json:"provisioningState,omitempty"`
	RecallStatus                                *ServerEndpointRecallStatus       `json:"recallStatus,omitempty"`
	ServerLocalPath                             *string                           `json:"serverLocalPath,omitempty"`
	ServerResourceId                            *string                           `json:"serverResourceId,omitempty"`
	SyncStatus                                  *ServerEndpointSyncStatus         `json:"syncStatus,omitempty"`
	TierFilesOlderThanDays                      *int64                            `json:"tierFilesOlderThanDays,omitempty"`
	VolumeFreeSpacePercent                      *int64                            `json:"volumeFreeSpacePercent,omitempty"`
}
