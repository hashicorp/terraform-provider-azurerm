package dedicatedhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedHostAllocatableVM struct {
	Count  *float64 `json:"count,omitempty"`
	VMSize *string  `json:"vmSize,omitempty"`
}
