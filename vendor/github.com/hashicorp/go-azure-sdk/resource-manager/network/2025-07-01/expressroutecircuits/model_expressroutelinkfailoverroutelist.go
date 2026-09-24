package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverRouteList struct {
	BeforeSimulation *[]ExpressRouteLinkFailoverRoute `json:"beforeSimulation,omitempty"`
	DuringSimulation *[]ExpressRouteLinkFailoverRoute `json:"duringSimulation,omitempty"`
}
