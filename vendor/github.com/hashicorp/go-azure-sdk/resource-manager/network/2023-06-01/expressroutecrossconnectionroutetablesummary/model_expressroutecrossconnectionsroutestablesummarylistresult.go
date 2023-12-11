package expressroutecrossconnectionroutetablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionsRoutesTableSummaryListResult struct {
	NextLink *string                                          `json:"nextLink,omitempty"`
	Value    *[]ExpressRouteCrossConnectionRoutesTableSummary `json:"value,omitempty"`
}
