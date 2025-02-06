package sapdiskconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskDetails struct {
	DiskTier                  *string  `json:"diskTier,omitempty"`
	IopsReadWrite             *int64   `json:"iopsReadWrite,omitempty"`
	MaximumSupportedDiskCount *int64   `json:"maximumSupportedDiskCount,omitempty"`
	MbpsReadWrite             *int64   `json:"mbpsReadWrite,omitempty"`
	MinimumSupportedDiskCount *int64   `json:"minimumSupportedDiskCount,omitempty"`
	SizeGB                    *int64   `json:"sizeGB,omitempty"`
	Sku                       *DiskSku `json:"sku,omitempty"`
}
