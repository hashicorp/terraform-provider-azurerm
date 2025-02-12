package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FleetUpdate struct {
	Identity   *ManagedServiceIdentityUpdate `json:"identity,omitempty"`
	Plan       *ResourcePlanUpdate           `json:"plan,omitempty"`
	Properties *FleetProperties              `json:"properties,omitempty"`
	Tags       *map[string]string            `json:"tags,omitempty"`
}
