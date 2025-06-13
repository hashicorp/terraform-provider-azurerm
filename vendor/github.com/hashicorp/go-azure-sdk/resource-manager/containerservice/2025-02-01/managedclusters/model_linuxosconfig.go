package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxOSConfig struct {
	SwapFileSizeMB             *int64        `json:"swapFileSizeMB,omitempty"`
	Sysctls                    *SysctlConfig `json:"sysctls,omitempty"`
	TransparentHugePageDefrag  *string       `json:"transparentHugePageDefrag,omitempty"`
	TransparentHugePageEnabled *string       `json:"transparentHugePageEnabled,omitempty"`
}
