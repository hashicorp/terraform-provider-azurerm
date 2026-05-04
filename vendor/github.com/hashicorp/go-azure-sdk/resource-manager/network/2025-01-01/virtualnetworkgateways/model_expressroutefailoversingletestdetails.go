package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteFailoverSingleTestDetails struct {
	EndTimeUtc                *string                               `json:"endTimeUtc,omitempty"`
	FailoverConnectionDetails *[]FailoverConnectionDetails          `json:"failoverConnectionDetails,omitempty"`
	NonRedundantRoutes        *[]string                             `json:"nonRedundantRoutes,omitempty"`
	PeeringLocation           *string                               `json:"peeringLocation,omitempty"`
	RedundantRoutes           *[]ExpressRouteFailoverRedundantRoute `json:"redundantRoutes,omitempty"`
	StartTimeUtc              *string                               `json:"startTimeUtc,omitempty"`
	Status                    *FailoverTestStatusForSingleTest      `json:"status,omitempty"`
	WasSimulationSuccessful   *bool                                 `json:"wasSimulationSuccessful,omitempty"`
}
