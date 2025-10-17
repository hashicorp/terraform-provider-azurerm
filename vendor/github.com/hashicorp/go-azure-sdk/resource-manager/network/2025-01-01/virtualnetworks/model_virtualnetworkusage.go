package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkUsage struct {
	CurrentValue *float64                 `json:"currentValue,omitempty"`
	Id           *string                  `json:"id,omitempty"`
	Limit        *float64                 `json:"limit,omitempty"`
	Name         *VirtualNetworkUsageName `json:"name,omitempty"`
	Unit         *string                  `json:"unit,omitempty"`
}
