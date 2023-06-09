package expressroutecrossconnectionroutetablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionRoutesTableSummary struct {
	Asn                     *int64  `json:"asn,omitempty"`
	Neighbor                *string `json:"neighbor,omitempty"`
	StateOrPrefixesReceived *string `json:"stateOrPrefixesReceived,omitempty"`
	UpDown                  *string `json:"upDown,omitempty"`
}
