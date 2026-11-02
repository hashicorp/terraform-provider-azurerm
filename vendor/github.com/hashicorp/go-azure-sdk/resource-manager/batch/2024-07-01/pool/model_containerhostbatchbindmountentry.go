package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerHostBatchBindMountEntry struct {
	IsReadOnly *bool                  `json:"isReadOnly,omitempty"`
	Source     *ContainerHostDataPath `json:"source,omitempty"`
}
