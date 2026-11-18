package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualHubRoute struct {
	AddressPrefixes  *[]string `json:"addressPrefixes,omitempty"`
	NextHopIPAddress *string   `json:"nextHopIpAddress,omitempty"`
}
