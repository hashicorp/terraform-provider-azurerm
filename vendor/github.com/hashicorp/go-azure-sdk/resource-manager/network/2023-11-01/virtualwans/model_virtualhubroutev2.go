package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualHubRouteV2 struct {
	DestinationType *string   `json:"destinationType,omitempty"`
	Destinations    *[]string `json:"destinations,omitempty"`
	NextHopType     *string   `json:"nextHopType,omitempty"`
	NextHops        *[]string `json:"nextHops,omitempty"`
}
