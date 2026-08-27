package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDisksToDetach struct {
	DetachOption *DiskDetachOptionTypes `json:"detachOption,omitempty"`
	DiskId       string                 `json:"diskId"`
}
