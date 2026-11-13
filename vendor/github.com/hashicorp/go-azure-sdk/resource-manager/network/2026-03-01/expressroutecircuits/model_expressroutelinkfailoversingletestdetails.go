package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverSingleTestDetails struct {
	BgpStatus               *[]ExpressRouteLinkFailoverTestBgpStatus `json:"bgpStatus,omitempty"`
	CircuitTestCategory     *MaintenanceTestCategory                 `json:"circuitTestCategory,omitempty"`
	EndTimeUtc              *string                                  `json:"endTimeUtc,omitempty"`
	IsSimulationVerified    *bool                                    `json:"isSimulationVerified,omitempty"`
	LinkType                *ExpressRouteFailoverLinkType            `json:"linkType,omitempty"`
	NonRedundantRoutes      *ExpressRouteLinkFailoverRouteList       `json:"nonRedundantRoutes,omitempty"`
	RedundantRoutes         *ExpressRouteLinkFailoverRouteList       `json:"redundantRoutes,omitempty"`
	StartTimeUtc            *string                                  `json:"startTimeUtc,omitempty"`
	Status                  *FailoverTestStatus                      `json:"status,omitempty"`
	WasSimulationSuccessful *bool                                    `json:"wasSimulationSuccessful,omitempty"`
}
