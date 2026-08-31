package expressroutecrossconnectionroutetable

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteCircuitRoutesTable struct {
	LocPrf  *string `json:"locPrf,omitempty"`
	Network *string `json:"network,omitempty"`
	NextHop *string `json:"nextHop,omitempty"`
	Path    *string `json:"path,omitempty"`
	Weight  *int64  `json:"weight,omitempty"`
}
