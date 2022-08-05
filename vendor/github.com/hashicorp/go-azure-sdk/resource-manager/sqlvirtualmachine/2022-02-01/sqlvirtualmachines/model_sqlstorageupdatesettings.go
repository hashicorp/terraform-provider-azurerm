package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlStorageUpdateSettings struct {
	DiskConfigurationType *DiskConfigurationType `json:"diskConfigurationType,omitempty"`
	DiskCount             *int64                 `json:"diskCount,omitempty"`
	StartingDeviceId      *int64                 `json:"startingDeviceId,omitempty"`
}
