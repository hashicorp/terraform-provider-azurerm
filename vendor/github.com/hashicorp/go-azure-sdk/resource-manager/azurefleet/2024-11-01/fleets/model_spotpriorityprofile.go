package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpotPriorityProfile struct {
	AllocationStrategy *SpotAllocationStrategy `json:"allocationStrategy,omitempty"`
	Capacity           *int64                  `json:"capacity,omitempty"`
	EvictionPolicy     *EvictionPolicy         `json:"evictionPolicy,omitempty"`
	Maintain           *bool                   `json:"maintain,omitempty"`
	MaxPricePerVM      *float64                `json:"maxPricePerVM,omitempty"`
	MinCapacity        *int64                  `json:"minCapacity,omitempty"`
}
