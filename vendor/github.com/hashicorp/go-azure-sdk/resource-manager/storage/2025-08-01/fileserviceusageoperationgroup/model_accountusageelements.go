package fileserviceusageoperationgroup

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountUsageElements struct {
	FileShareCount                *int64 `json:"fileShareCount,omitempty"`
	ProvisionedBandwidthMiBPerSec *int64 `json:"provisionedBandwidthMiBPerSec,omitempty"`
	ProvisionedIOPS               *int64 `json:"provisionedIOPS,omitempty"`
	ProvisionedStorageGiB         *int64 `json:"provisionedStorageGiB,omitempty"`
}
