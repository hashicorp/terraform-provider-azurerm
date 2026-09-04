package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Disk struct {
	DiskType         *string `json:"diskType,omitempty"`
	GeneratedId      *string `json:"generatedId,omitempty"`
	Id               *string `json:"id,omitempty"`
	MaxSizeInBytes   *int64  `json:"maxSizeInBytes,omitempty"`
	Name             *string `json:"name,omitempty"`
	Path             *string `json:"path,omitempty"`
	UsedSpaceInBytes *int64  `json:"usedSpaceInBytes,omitempty"`
}
