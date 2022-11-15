package capacities

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DedicatedCapacityMutableProperties struct {
	Administration *DedicatedCapacityAdministrators `json:"administration,omitempty"`
	FriendlyName   *string                          `json:"friendlyName,omitempty"`
	Mode           *Mode                            `json:"mode,omitempty"`
	TenantId       *string                          `json:"tenantId,omitempty"`
}
