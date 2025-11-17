package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMSSDataDisk struct {
	DiskLetter string   `json:"diskLetter"`
	DiskSizeGB int64    `json:"diskSizeGB"`
	DiskType   DiskType `json:"diskType"`
	Lun        int64    `json:"lun"`
}
