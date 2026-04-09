package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccountSkuConversionStatus struct {
	EndTime             *string              `json:"endTime,omitempty"`
	SkuConversionStatus *SkuConversionStatus `json:"skuConversionStatus,omitempty"`
	StartTime           *string              `json:"startTime,omitempty"`
	TargetSkuName       *SkuName             `json:"targetSkuName,omitempty"`
}
