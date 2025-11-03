package routingrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRuleNextHop struct {
	NextHopAddress *string                `json:"nextHopAddress,omitempty"`
	NextHopType    RoutingRuleNextHopType `json:"nextHopType"`
}
