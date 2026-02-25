package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualHubEffectiveRoute struct {
	AddressPrefixes *[]string `json:"addressPrefixes,omitempty"`
	AsPath          *string   `json:"asPath,omitempty"`
	NextHopType     *string   `json:"nextHopType,omitempty"`
	NextHops        *[]string `json:"nextHops,omitempty"`
	RouteOrigin     *string   `json:"routeOrigin,omitempty"`
}
