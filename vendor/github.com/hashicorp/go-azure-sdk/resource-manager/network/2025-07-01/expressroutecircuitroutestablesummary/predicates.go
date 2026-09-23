package expressroutecircuitroutestablesummary

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitRoutesTableSummaryOperationPredicate struct {
	As          *int64
	Neighbor    *string
	StatePfxRcd *string
	UpDown      *string
	V           *int64
}

func (p ExpressRouteCircuitRoutesTableSummaryOperationPredicate) Matches(input ExpressRouteCircuitRoutesTableSummary) bool {

	if p.As != nil && (input.As == nil || *p.As != *input.As) {
		return false
	}

	if p.Neighbor != nil && (input.Neighbor == nil || *p.Neighbor != *input.Neighbor) {
		return false
	}

	if p.StatePfxRcd != nil && (input.StatePfxRcd == nil || *p.StatePfxRcd != *input.StatePfxRcd) {
		return false
	}

	if p.UpDown != nil && (input.UpDown == nil || *p.UpDown != *input.UpDown) {
		return false
	}

	if p.V != nil && (input.V == nil || *p.V != *input.V) {
		return false
	}

	return true
}
