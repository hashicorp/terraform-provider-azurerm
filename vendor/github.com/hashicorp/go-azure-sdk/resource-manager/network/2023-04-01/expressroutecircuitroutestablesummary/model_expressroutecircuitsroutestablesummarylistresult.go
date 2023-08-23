package expressroutecircuitroutestablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitsRoutesTableSummaryListResult struct {
	NextLink *string                                  `json:"nextLink,omitempty"`
	Value    *[]ExpressRouteCircuitRoutesTableSummary `json:"value,omitempty"`
}
