package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointCreateParametersProperties struct {
	CloudTiering                 *FeatureStatus         `json:"cloudTiering,omitempty"`
	FriendlyName                 *string                `json:"friendlyName,omitempty"`
	InitialDownloadPolicy        *InitialDownloadPolicy `json:"initialDownloadPolicy,omitempty"`
	LocalCacheMode               *LocalCacheMode        `json:"localCacheMode,omitempty"`
	OfflineDataTransfer          *FeatureStatus         `json:"offlineDataTransfer,omitempty"`
	OfflineDataTransferShareName *string                `json:"offlineDataTransferShareName,omitempty"`
	ServerLocalPath              *string                `json:"serverLocalPath,omitempty"`
	ServerResourceId             *string                `json:"serverResourceId,omitempty"`
	TierFilesOlderThanDays       *int64                 `json:"tierFilesOlderThanDays,omitempty"`
	VolumeFreeSpacePercent       *int64                 `json:"volumeFreeSpacePercent,omitempty"`
}
