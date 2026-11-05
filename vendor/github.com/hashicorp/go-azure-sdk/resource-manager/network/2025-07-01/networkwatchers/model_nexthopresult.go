package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NextHopResult struct {
	NextHopIPAddress *string      `json:"nextHopIpAddress,omitempty"`
	NextHopType      *NextHopType `json:"nextHopType,omitempty"`
	RouteTableId     *string      `json:"routeTableId,omitempty"`
}
