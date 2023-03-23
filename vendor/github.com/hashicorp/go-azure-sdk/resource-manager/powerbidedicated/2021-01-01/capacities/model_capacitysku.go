package capacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CapacitySku struct {
	Capacity *int64           `json:"capacity,omitempty"`
	Name     string           `json:"name"`
	Tier     *CapacitySkuTier `json:"tier,omitempty"`
}
