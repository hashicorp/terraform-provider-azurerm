package sapdiskconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskVolumeConfiguration struct {
	Count  *int64   `json:"count,omitempty"`
	SizeGB *int64   `json:"sizeGB,omitempty"`
	Sku    *DiskSku `json:"sku,omitempty"`
}
