package expressroutecrossconnectionarptable

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitArpTableOperationPredicate struct {
	Age        *int64
	IPAddress  *string
	Interface  *string
	MacAddress *string
}

func (p ExpressRouteCircuitArpTableOperationPredicate) Matches(input ExpressRouteCircuitArpTable) bool {

	if p.Age != nil && (input.Age == nil || *p.Age != *input.Age) {
		return false
	}

	if p.IPAddress != nil && (input.IPAddress == nil || *p.IPAddress != *input.IPAddress) {
		return false
	}

	if p.Interface != nil && (input.Interface == nil || *p.Interface != *input.Interface) {
		return false
	}

	if p.MacAddress != nil && (input.MacAddress == nil || *p.MacAddress != *input.MacAddress) {
		return false
	}

	return true
}
