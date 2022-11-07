package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SnapshotSku struct {
	Name *SnapshotStorageAccountTypes `json:"name,omitempty"`
	Tier *string                      `json:"tier,omitempty"`
}
