package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HardwareProfile struct {
	VMSize           *VirtualMachineSizeTypes `json:"vmSize,omitempty"`
	VMSizeProperties *VMSizeProperties        `json:"vmSizeProperties,omitempty"`
}
