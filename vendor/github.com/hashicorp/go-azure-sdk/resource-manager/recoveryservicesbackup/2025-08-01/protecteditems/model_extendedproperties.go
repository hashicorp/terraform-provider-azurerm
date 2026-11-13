package protecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtendedProperties struct {
	DiskExclusionProperties *DiskExclusionProperties `json:"diskExclusionProperties,omitempty"`
	LinuxVMApplicationName  *string                  `json:"linuxVmApplicationName,omitempty"`
}
