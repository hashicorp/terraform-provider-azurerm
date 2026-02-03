package expressroutecrossconnectionroutetablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCrossConnectionRoutesTableSummaryOperationPredicate struct {
	Asn                     *int64
	Neighbor                *string
	StateOrPrefixesReceived *string
	UpDown                  *string
}

func (p ExpressRouteCrossConnectionRoutesTableSummaryOperationPredicate) Matches(input ExpressRouteCrossConnectionRoutesTableSummary) bool {

	if p.Asn != nil && (input.Asn == nil || *p.Asn != *input.Asn) {
		return false
	}

	if p.Neighbor != nil && (input.Neighbor == nil || *p.Neighbor != *input.Neighbor) {
		return false
	}

	if p.StateOrPrefixesReceived != nil && (input.StateOrPrefixesReceived == nil || *p.StateOrPrefixesReceived != *input.StateOrPrefixesReceived) {
		return false
	}

	if p.UpDown != nil && (input.UpDown == nil || *p.UpDown != *input.UpDown) {
		return false
	}

	return true
}
