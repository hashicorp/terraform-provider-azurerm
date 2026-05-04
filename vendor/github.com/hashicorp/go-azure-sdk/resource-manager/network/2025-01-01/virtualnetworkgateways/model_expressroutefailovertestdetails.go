package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteFailoverTestDetails struct {
	Circuits        *[]ExpressRouteFailoverCircuitResourceDetails    `json:"circuits,omitempty"`
	Connections     *[]ExpressRouteFailoverConnectionResourceDetails `json:"connections,omitempty"`
	EndTime         *string                                          `json:"endTime,omitempty"`
	Issues          *[]string                                        `json:"issues,omitempty"`
	PeeringLocation *string                                          `json:"peeringLocation,omitempty"`
	StartTime       *string                                          `json:"startTime,omitempty"`
	Status          *FailoverTestStatus                              `json:"status,omitempty"`
	TestGuid        *string                                          `json:"testGuid,omitempty"`
	TestType        *FailoverTestType                                `json:"testType,omitempty"`
}
