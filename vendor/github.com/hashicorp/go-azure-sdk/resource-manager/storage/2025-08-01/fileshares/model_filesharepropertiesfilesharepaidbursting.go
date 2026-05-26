package fileshares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FileSharePropertiesFileSharePaidBursting struct {
	PaidBurstingEnabled           *bool  `json:"paidBurstingEnabled,omitempty"`
	PaidBurstingMaxBandwidthMibps *int64 `json:"paidBurstingMaxBandwidthMibps,omitempty"`
	PaidBurstingMaxIops           *int64 `json:"paidBurstingMaxIops,omitempty"`
}
