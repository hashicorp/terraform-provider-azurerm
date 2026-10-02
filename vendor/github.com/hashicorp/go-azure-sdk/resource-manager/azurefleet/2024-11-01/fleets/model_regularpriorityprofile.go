package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegularPriorityProfile struct {
	AllocationStrategy *RegularPriorityAllocationStrategy `json:"allocationStrategy,omitempty"`
	Capacity           *int64                             `json:"capacity,omitempty"`
	MinCapacity        *int64                             `json:"minCapacity,omitempty"`
}
