package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Placement struct {
	ExcludeZones        *[]string                `json:"excludeZones,omitempty"`
	IncludeZones        *[]string                `json:"includeZones,omitempty"`
	ZonePlacementPolicy *ZonePlacementPolicyType `json:"zonePlacementPolicy,omitempty"`
}
