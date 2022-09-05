package capacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedCapacityUpdateParameters struct {
	Properties *DedicatedCapacityMutableProperties `json:"properties,omitempty"`
	Sku        *CapacitySku                        `json:"sku,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
}
