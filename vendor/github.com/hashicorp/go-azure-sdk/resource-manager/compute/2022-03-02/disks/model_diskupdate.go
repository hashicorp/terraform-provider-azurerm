package disks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiskUpdate struct {
	Properties *DiskUpdateProperties `json:"properties"`
	Sku        *DiskSku              `json:"sku"`
	Tags       *map[string]string    `json:"tags,omitempty"`
}
