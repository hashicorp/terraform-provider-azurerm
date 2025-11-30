package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayRouteSet struct {
	Details   *map[string][]RouteSourceDetails `json:"details,omitempty"`
	Locations *[]string                        `json:"locations,omitempty"`
	Name      *string                          `json:"name,omitempty"`
}
