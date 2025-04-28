package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTargetSpaceAllocation struct {
	AllocationPercentage *int64  `json:"allocationPercentage,omitempty"`
	Name                 *string `json:"name,omitempty"`
}
