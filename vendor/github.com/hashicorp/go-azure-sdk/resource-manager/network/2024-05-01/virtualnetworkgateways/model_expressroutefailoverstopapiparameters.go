package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteFailoverStopApiParameters struct {
	Details                 *[]FailoverConnectionDetails `json:"details,omitempty"`
	PeeringLocation         *string                      `json:"peeringLocation,omitempty"`
	WasSimulationSuccessful *bool                        `json:"wasSimulationSuccessful,omitempty"`
}
