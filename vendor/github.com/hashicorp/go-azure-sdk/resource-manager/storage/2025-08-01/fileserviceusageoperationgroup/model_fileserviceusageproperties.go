package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileServiceUsageProperties struct {
	BurstingConstants        *BurstingConstants        `json:"burstingConstants,omitempty"`
	FileShareLimits          *FileShareLimits          `json:"fileShareLimits,omitempty"`
	FileShareRecommendations *FileShareRecommendations `json:"fileShareRecommendations,omitempty"`
	StorageAccountLimits     *AccountLimits            `json:"storageAccountLimits,omitempty"`
	StorageAccountUsage      *AccountUsage             `json:"storageAccountUsage,omitempty"`
}
