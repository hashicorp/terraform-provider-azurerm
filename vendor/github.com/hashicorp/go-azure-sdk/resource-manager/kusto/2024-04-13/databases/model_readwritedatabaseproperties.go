package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReadWriteDatabaseProperties struct {
	HotCachePeriod     *string             `json:"hotCachePeriod,omitempty"`
	IsFollowed         *bool               `json:"isFollowed,omitempty"`
	KeyVaultProperties *KeyVaultProperties `json:"keyVaultProperties,omitempty"`
	ProvisioningState  *ProvisioningState  `json:"provisioningState,omitempty"`
	SoftDeletePeriod   *string             `json:"softDeletePeriod,omitempty"`
	Statistics         *DatabaseStatistics `json:"statistics,omitempty"`
	SuspensionDetails  *SuspensionDetails  `json:"suspensionDetails,omitempty"`
}
