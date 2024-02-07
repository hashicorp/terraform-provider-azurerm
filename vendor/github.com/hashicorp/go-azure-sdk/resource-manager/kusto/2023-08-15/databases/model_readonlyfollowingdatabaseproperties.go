package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReadOnlyFollowingDatabaseProperties struct {
	AttachedDatabaseConfigurationName *string                      `json:"attachedDatabaseConfigurationName,omitempty"`
	DatabaseShareOrigin               *DatabaseShareOrigin         `json:"databaseShareOrigin,omitempty"`
	HotCachePeriod                    *string                      `json:"hotCachePeriod,omitempty"`
	LeaderClusterResourceId           *string                      `json:"leaderClusterResourceId,omitempty"`
	OriginalDatabaseName              *string                      `json:"originalDatabaseName,omitempty"`
	PrincipalsModificationKind        *PrincipalsModificationKind  `json:"principalsModificationKind,omitempty"`
	ProvisioningState                 *ProvisioningState           `json:"provisioningState,omitempty"`
	SoftDeletePeriod                  *string                      `json:"softDeletePeriod,omitempty"`
	Statistics                        *DatabaseStatistics          `json:"statistics,omitempty"`
	SuspensionDetails                 *SuspensionDetails           `json:"suspensionDetails,omitempty"`
	TableLevelSharingProperties       *TableLevelSharingProperties `json:"tableLevelSharingProperties,omitempty"`
}
