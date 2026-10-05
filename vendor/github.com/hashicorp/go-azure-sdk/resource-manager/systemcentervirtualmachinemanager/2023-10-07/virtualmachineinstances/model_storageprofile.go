package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageProfile struct {
	Disks *[]VirtualDisk `json:"disks,omitempty"`
}
