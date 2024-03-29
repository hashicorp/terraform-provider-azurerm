package logicalnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutePropertiesFormat struct {
	AddressPrefix    *string `json:"addressPrefix,omitempty"`
	NextHopIPAddress *string `json:"nextHopIpAddress,omitempty"`
}
