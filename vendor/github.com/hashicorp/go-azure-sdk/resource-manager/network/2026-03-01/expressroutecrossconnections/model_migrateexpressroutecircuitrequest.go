package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateExpressRouteCircuitRequest struct {
	PortId                *string        `json:"portId,omitempty"`
	TargetPeeringLocation *string        `json:"targetPeeringLocation,omitempty"`
	TargetPortMapping     *[]PortMapping `json:"targetPortMapping,omitempty"`
}
