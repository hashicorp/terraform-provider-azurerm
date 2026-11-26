package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverAllTestsDetails struct {
	BgpStatus               *[]ExpressRouteLinkFailoverTestBgpStatus `json:"bgpStatus,omitempty"`
	CircuitTestCategory     *MaintenanceTestCategory                 `json:"circuitTestCategory,omitempty"`
	EndTime                 *string                                  `json:"endTime,omitempty"`
	Issues                  *[]string                                `json:"issues,omitempty"`
	LinkType                *ExpressRouteFailoverLinkType            `json:"linkType,omitempty"`
	StartTime               *string                                  `json:"startTime,omitempty"`
	Status                  *FailoverTestStatus                      `json:"status,omitempty"`
	TestGuid                *string                                  `json:"testGuid,omitempty"`
	TestType                *FailoverTestType                        `json:"testType,omitempty"`
	WasSimulationSuccessful *bool                                    `json:"wasSimulationSuccessful,omitempty"`
}
