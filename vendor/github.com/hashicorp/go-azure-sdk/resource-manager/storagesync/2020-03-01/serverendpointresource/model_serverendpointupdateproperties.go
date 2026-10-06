package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointUpdateProperties struct {
	CloudTiering                 *FeatureStatus  `json:"cloudTiering,omitempty"`
	LocalCacheMode               *LocalCacheMode `json:"localCacheMode,omitempty"`
	OfflineDataTransfer          *FeatureStatus  `json:"offlineDataTransfer,omitempty"`
	OfflineDataTransferShareName *string         `json:"offlineDataTransferShareName,omitempty"`
	TierFilesOlderThanDays       *int64          `json:"tierFilesOlderThanDays,omitempty"`
	VolumeFreeSpacePercent       *int64          `json:"volumeFreeSpacePercent,omitempty"`
}
