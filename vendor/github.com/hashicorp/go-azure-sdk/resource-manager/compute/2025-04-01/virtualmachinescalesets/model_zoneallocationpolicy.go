package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZoneAllocationPolicy struct {
	MaxInstancePercentPerZonePolicy *MaxInstancePercentPerZonePolicy `json:"maxInstancePercentPerZonePolicy,omitempty"`
	MaxZoneCount                    *int64                           `json:"maxZoneCount,omitempty"`
}
