package packetcorecontrolplane

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InterfaceProperties struct {
	IPv4Address *string `json:"ipv4Address,omitempty"`
	IPv4Gateway *string `json:"ipv4Gateway,omitempty"`
	IPv4Subnet  *string `json:"ipv4Subnet,omitempty"`
	Name        *string `json:"name,omitempty"`
}
