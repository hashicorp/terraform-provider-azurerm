package labs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataDiskProperties struct {
	AttachNewDataDiskOptions *AttachNewDataDiskOptions `json:"attachNewDataDiskOptions,omitempty"`
	ExistingLabDiskId        *string                   `json:"existingLabDiskId,omitempty"`
	HostCaching              *HostCachingOptions       `json:"hostCaching,omitempty"`
}
