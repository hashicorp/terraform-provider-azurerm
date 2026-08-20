package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverStopApiParameters struct {
	CircuitTestCategory     *string `json:"circuitTestCategory,omitempty"`
	IsVerified              *bool   `json:"isVerified,omitempty"`
	LinkType                *string `json:"linkType,omitempty"`
	WasSimulationSuccessful *bool   `json:"wasSimulationSuccessful,omitempty"`
}
