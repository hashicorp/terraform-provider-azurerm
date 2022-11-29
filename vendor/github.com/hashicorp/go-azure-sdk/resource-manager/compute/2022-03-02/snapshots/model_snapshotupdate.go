package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotUpdate struct {
	Properties *SnapshotUpdateProperties `json:"properties"`
	Sku        *SnapshotSku              `json:"sku"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}
