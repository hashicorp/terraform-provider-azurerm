package expressroutecircuitroutestable

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitRoutesTableOperationPredicate struct {
	LocPrf  *string
	Network *string
	NextHop *string
	Path    *string
	Weight  *int64
}

func (p ExpressRouteCircuitRoutesTableOperationPredicate) Matches(input ExpressRouteCircuitRoutesTable) bool {

	if p.LocPrf != nil && (input.LocPrf == nil || *p.LocPrf != *input.LocPrf) {
		return false
	}

	if p.Network != nil && (input.Network == nil || *p.Network != *input.Network) {
		return false
	}

	if p.NextHop != nil && (input.NextHop == nil || *p.NextHop != *input.NextHop) {
		return false
	}

	if p.Path != nil && (input.Path == nil || *p.Path != *input.Path) {
		return false
	}

	if p.Weight != nil && (input.Weight == nil || *p.Weight != *input.Weight) {
		return false
	}

	return true
}
