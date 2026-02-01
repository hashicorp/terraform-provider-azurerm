package vpngateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticRoute struct {
	AddressPrefixes  *[]string `json:"addressPrefixes,omitempty"`
	Name             *string   `json:"name,omitempty"`
	NextHopIPAddress *string   `json:"nextHopIpAddress,omitempty"`
}
